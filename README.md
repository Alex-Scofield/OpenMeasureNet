# OpenMeasureNet
**OpenMeasureNet** is a free, open source system for building and monitoring networks of measurement devices. It aims to provide communities with the ability to share real-time measurement data among their members and, if desired, to the public, easily and securely. 

**OpenMeasureNet** takes care of the boilerplate that connects the measurements taken by your systems to a database, and that database to your measurement devices' owners. That way, you can build your own network for your community by modifying a single configuration file. Of course, **OpenMeasureNet** is highly customizable, but you should be able to get your network up and running with very little work. In particular, you may, if you wish, replace any of the components mentioned below to alter the system's functionalities, as long as the essential interfaces use the correct formats.

An **OpenMeasureNet** project is made out of the following components:
- *OMN Nodes*: An node is a device that measures and transmits data to the network. This would typically be an IoT device with sensors, but it could also be a computer sharing its ressource usage or even a program that transmits manual entries from a user. to the network Each node is owned by a user and each user can own any number of nodes. Users can choose whether they want their data to be shared with the rest of the network.
- *OMN Network*: The strength of **OpenMeasureNet** is that it allows for users to share the data they have captured with their entire community (or even, with multiple different communities), if desired, in real time. The Network component is essentially a broker that collects all of the data, a program which processes it and stores it in a database and a RESTful API for users to interact with this database and to fetch customizable dashbords to view the data. The default database is an implementation of the Entity-Relationship diagram below. However, *you* are in control of *your* data, and so can configure **OpenMeasureNet** to work with your existing database. 

Given how varied projects may be, it is important that they be easily extendable, for instance by adding API endpoints. In order to not bloat up the base system, which is intended to be as simple as possible, we include some additional features that you can opt in on when creating your projects. 
- *OMN Frontend*: You are free to use the backend features whilst building your own frontend. This allows you, for instance, to easily integrate this project with your existing infrastructure. You may, however, choose to use a default *OMN Frontend* to get started more easily. 
- *OMN Local*: Users may find it useful to monitor their devices without connecting to the network. Through *OMN Local*, a mobile and desktop application, nodes can be monitored directly, through a LAN or via Bluetooth. It is also possible to use this *plugin* to send data to the network through a local device (this can be useful, for instance, to transfer data in bulk).
- *OMN Forecast*: This is a tool that implements a LSTM model for forecasting measurements for each individual node in the network. It includes a training script, to be ran on a schedule and a prediction script, together with an API to fetch the results.

## OMN Network
The following entity-relationship diagram presents the fundamental entities that are involved in an *OMN Network*. These can, of course, be extended, and more entities can be added, but any implementation should include at least the following.

![docs/media/er-diagram.jpg](ER Diagram)

The *OMN Network* can be separated into three components:
1. *OMN Transfer*: This is the component that takes care of *transferring* the data from the users to the database. This can be achieved via an *mqtt* broker, for realtime monitoring, or a *mass transfer* component, where large quantities of measurements are transferred all at once, which can be useful when some devices are not connected to the internet, or when users may wish to not share their data in real time for privacy concerns. Other transfer systems, such as ones using LoRa or even Bluetooth could also be implemented. This component can be reused, for instance, for transferring data locally.
2. *OMN Backend*: This is a backend that users interact with through a RESTful API that defines how users log into the system and fetch the data they require. By construction it can be extended easily.
3. *OMN Update Manager*: Updates may need to be distributed to the measurement devices. This is done through the Update Manager component.

The following sequence diagram shows the interactions that happen when a user first turns on a node for realtime transfer, which sends two messages, and then logs into the backend and fetches a dashboard to view its data.

![docs/media/sequence-diagram.jpg](Sequence Diagram)

## OMN Payload
To transfer data from a node to the network, **OpenMeasureNet** follows the simple format below:

![docs/media/protocol-diagram.jpg](Protocol Diagram)
This format was chosen for the following reasons:
1. It is lightweight, and thus suitable for devices that have low network connectivity.
2. Each individual message transferred is semantically complete (as it is a measured quantity, at a given place, at a given time), and does not contain personal data.
3. It allows for authorization and authentication to take place in the application layer (thus preventing security vulnerabilities such as spoofing).
NOTE: TIMESTAMP IS ACTUALLY FLOAT64.
There is a standard list of *quantity* identifiers that we encourage users to respect and to consult when developing their own node software. This is because, that way, different measurement systems can operate in different networks seamlessly, thus reducing vendor lock-in and allowing maximum flexibility with respect to the devices that can work with each network. 

*Quantity ids* (or *type*, above) are given by a 1 byte unsigned integer. Quantities from 0 to 127 are reserved for this list, and those from 128 to 255 are free to represent any other quantities, as defined by each network's developers. 

## Architecture
The following diagram illustrates the system's architecture. In white, the core features, in red the *OMN Local* plugin, in green the *OMN Forecast* plugin and in purple the *OMN Frontend* plugin. 

![media/docs/component-diagram.jpg](Component Diagram)
## How to use **OpenMeasureNet**
1. Clone this repository
2. Clone any desired plugins
3. Develop or install **OpenMeasureNet** compliant firmware for your nodes
4. Fill in the *.env* file with the relevant configuration
5. Configure encryption (for instance SSL / TLS)
6. Deploy the base and any desired plugins via docker compose
