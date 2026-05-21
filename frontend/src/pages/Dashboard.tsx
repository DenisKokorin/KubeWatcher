import React, { useEffect, useState } from 'react'
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  CircularProgress,
  Alert,
} from '@mui/material'
import { useSEO } from '../hooks/useSEO'
import {
  Storage as NodesIcon,
  ViewList as PodsIcon,
  Apps as DeploymentsIcon,
  SettingsEthernet as ServicesIcon,
} from '@mui/icons-material'
import { monitorService, Cluster, NodeSummary, PodSummary, DeploymentSummary, ServiceSummary } from '../services/monitorService'

const Dashboard: React.FC = () => {
  const [cluster, setCluster] = useState<Cluster | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useSEO({
    title: 'Dashboard',
    description: 'KubeWatcher dashboard provides an overview of Kubernetes cluster health and team review status.',
    path: '/dashboard',
    noindex: true,
  })

  useEffect(() => {
    fetchClusterData()
  }, [])

  const fetchClusterData = async () => {
    try {
      const data = await monitorService.getClusterOverview()
      setCluster(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch cluster data')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>
  }

  if (!cluster) {
    return <Alert severity="info">No cluster data available</Alert>
  }

  const stats = [
    {
      title: 'Nodes',
      icon: <NodesIcon sx={{ fontSize: 40, color: '#2196f3' }} />,
      data: cluster.nodes,
      color: '#2196f3',
      type: 'nodes' as const,
    },
    {
      title: 'Pods',
      icon: <PodsIcon sx={{ fontSize: 40, color: '#4caf50' }} />,
      data: cluster.pods,
      color: '#4caf50',
      type: 'pods' as const,
    },
    {
      title: 'Deployments',
      icon: <DeploymentsIcon sx={{ fontSize: 40, color: '#ff9800' }} />,
      data: cluster.deployments,
      color: '#ff9800',
      type: 'deployments' as const,
    },
    {
      title: 'Services',
      icon: <ServicesIcon sx={{ fontSize: 40, color: '#9c27b0' }} />,
      data: cluster.services,
      color: '#9c27b0',
      type: 'services' as const,
    },
  ]

  return (
    <Box>
      <Typography variant="h4" gutterBottom sx={{ color: 'white', mb: 3 }}>
        Cluster Overview
      </Typography>

      <Grid container spacing={3}>
        {stats.map((stat) => (
          <Grid item xs={12} sm={6} md={3} key={stat.title}>
            <Card
              sx={{
                backgroundColor: '#1e1e1e',
                color: 'white',
                border: `1px solid ${stat.color}20`,
                '&:hover': {
                  borderColor: stat.color,
                  boxShadow: `0 0 10px ${stat.color}40`,
                },
              }}
            >
              <CardContent>
                <Box display="flex" alignItems="center" mb={2}>
                  {stat.icon}
                  <Typography variant="h6" sx={{ ml: 1 }}>
                    {stat.title}
                  </Typography>
                </Box>

                <Typography variant="h4" sx={{ mb: 1, color: stat.color }}>
                  {stat.data.total}
                </Typography>

                {stat.type === 'nodes' && (
                  <Box>
                    <Typography variant="body2" sx={{ color: '#4caf50' }}>
                      Ready: {(stat.data as NodeSummary).ready}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#f44336' }}>
                      Not Ready: {(stat.data as NodeSummary).not_ready}
                    </Typography>
                  </Box>
                )}

                {stat.type === 'pods' && (
                  <Box>
                    <Typography variant="body2" sx={{ color: '#4caf50' }}>
                      Running: {(stat.data as PodSummary).running}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#ff9800' }}>
                      Pending: {(stat.data as PodSummary).pending}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#f44336' }}>
                      Failed: {(stat.data as PodSummary).failed}
                    </Typography>
                  </Box>
                )}

                {stat.type === 'deployments' && (
                  <Box>
                    <Typography variant="body2" sx={{ color: '#4caf50' }}>
                      Ready: {(stat.data as DeploymentSummary).ready}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#f44336' }}>
                      Not Ready: {(stat.data as DeploymentSummary).not_ready}
                    </Typography>
                  </Box>
                )}

                {stat.type === 'services' && (
                  <Box>
                    <Typography variant="body2" sx={{ color: '#2196f3' }}>
                      ClusterIP: {(stat.data as ServiceSummary).cluster_ip}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#ff9800' }}>
                      LoadBalancer: {(stat.data as ServiceSummary).load_balancer}
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#9c27b0' }}>
                      NodePort: {(stat.data as ServiceSummary).node_port}
                    </Typography>
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Box mt={4}>
        <Typography variant="body2" sx={{ color: 'gray' }}>
          Last updated: {new Date(cluster.timestamp * 1000).toLocaleString()}
        </Typography>
      </Box>
    </Box>
  )
}

export default Dashboard