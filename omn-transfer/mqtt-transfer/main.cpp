#include <iostream>
#include <string>
#include <thread>
#include <chrono>
#include <cstdint>
#include <cstring>
#include <cstdlib>
#include <arpa/inet.h>
#include <endian.h>
#include <mqtt/client.h>
#include <pqxx/pqxx>

struct Observation {
    uint8_t quantity_id;
    float value;
    double timestamp;
    float longitude;
    float latitude;

    static Observation from_bytes(const std::string& raw) {
        if (raw.size() < 21)
            throw std::runtime_error(
                "bad payload (" + std::to_string(raw.size()) + " bytes, need 21)");

        Observation obs{};
        obs.quantity_id = raw[0];
        obs.value       = net_to_float(raw, 1);
        obs.timestamp   = net_to_double(raw, 5);
        obs.longitude   = net_to_float(raw, 13);
        obs.latitude    = net_to_float(raw, 17);
        return obs;
    }

private:
    static float net_to_float(const std::string& raw, size_t offset) {
        uint32_t network;
        std::memcpy(&network, &raw[offset], 4);
        network = ntohl(network);
        float host;
        std::memcpy(&host, &network, 4);
        return host;
    }

    static double net_to_double(const std::string& raw, size_t offset) {
        uint64_t network;
        std::memcpy(&network, &raw[offset], 8);
        network = be64toh(network);
        double host;
        std::memcpy(&host, &network, 8);
        return host;
    }
};

int main() {
    std::string mqtt_host = std::getenv("MQTT_HOST");
    std::string mqtt_client_id = std::getenv("MQTT_ARCHIVER_CLIENT_ID");
    std::string pg_host = std::getenv("POSTGRES_HOST");
    std::string pg_user = std::getenv("POSTGRES_USER");
    std::string pg_pass = std::getenv("POSTGRES_PASSWORD");
    std::string pg_db = std::getenv("POSTGRES_DB");
    std::string archiver_user = std::getenv("MQTT_ARCHIVER_USER");
    std::string archiver_pass = std::getenv("MQTT_ARCHIVER_PASSWORD");

    std::string pg_conn = "host=" + pg_host + " user=" + pg_user
        + " password=" + pg_pass + " dbname=" + pg_db;

    pqxx::connection pg(pg_conn);
    std::cout << "Connected to PostgreSQL" << std::endl;

    auto mqtt_opts = mqtt::connect_options_builder()
        .user_name(archiver_user)
        .password(archiver_pass)
        .clean_session(true)
        .finalize();

    while (true) {
        try {
            mqtt::client mqtt(mqtt_host, mqtt_client_id);
            mqtt.connect(mqtt_opts);
            mqtt.subscribe("o/#", 0);
            std::cout << "Connected to MQTT, subscribed to o/#" << std::endl;

            while (true) {
                try {
                    auto msg = mqtt.consume_message();
                    if (!msg) continue;

                    auto slash = msg->get_topic().find('/');
                    if (slash == std::string::npos) continue;
                    int boitier_id = std::stoi(
                        msg->get_topic().substr(slash + 1));

                    Observation obs = Observation::from_bytes(
                        msg->get_payload());

                    pqxx::work txn(pg);
                    txn.exec_params(
                        "INSERT INTO observations "
                        "(time, boitier_id, quantity, value, location) "
                        "VALUES (to_timestamp($1), $2, $3, $4, "
                        "ST_SetSRID(ST_MakePoint($5, $6), 4326)) "
                        "ON CONFLICT ON CONSTRAINT observations_pkey "
                        "DO NOTHING",
                        obs.timestamp, boitier_id, static_cast<int>(obs.quantity_id),
                        obs.value, obs.longitude, obs.latitude);
                    txn.commit();
                } catch (const std::exception& e) {
                    std::cerr << "Skipping bad message: "
                              << e.what() << std::endl;
                }
            }
        } catch (const std::exception& e) {
            std::cerr << "Connection lost (" << e.what()
                      << "), reconnecting in 3s..." << std::endl;
            std::this_thread::sleep_for(std::chrono::seconds(3));
        }
    }

    return 0;
}
