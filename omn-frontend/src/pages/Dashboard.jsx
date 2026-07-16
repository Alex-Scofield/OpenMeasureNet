import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { getNodes, getMap, logout } from '../api'

export default function Dashboard() {
  const [nodes, setNodes] = useState([])
  const [mapUrl, setMapUrl] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    async function load() {
      const nodesData = await getNodes()
      if (Array.isArray(nodesData)) setNodes(nodesData)

      const mapData = await getMap()
      if (mapData.dashboard_url) setMapUrl(mapData.dashboard_url)
    }
    load()
  }, [])

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>OpenMeasureNet</h1>
        <button onClick={handleLogout}>Log out</button>
      </header>

      <section className="dashboard-section">
        <h2>Map</h2>
        {mapUrl ? (
          <div className="iframe-container map-container">
            <iframe src={mapUrl} title="All Nodes Map" />
          </div>
        ) : (
          <p>Loading...</p>
        )}
      </section>

      <section className="dashboard-section">
        <h2>My Nodes</h2>
        {nodes.length === 0 ? (
          <p>No nodes yet.</p>
        ) : (
          <div className="nodes-grid">
            {nodes.map(node => (
              <div key={node.id} className="node-card">
                <h3>Node {node.id}</h3>
                <div className="iframe-container">
                  <iframe src={node.dashboard_url} title={`Node ${node.id}`} />
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
