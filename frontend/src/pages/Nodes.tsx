import React, { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  CircularProgress,
  Alert,
  Chip,
} from '@mui/material'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import { monitorService, Node } from '../services/monitorService'
import { useSEO } from '../hooks/useSEO'

const Nodes: React.FC = () => {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useSEO({
    title: 'Cluster Nodes',
    description: 'Monitor Kubernetes node status and health metrics inside KubeWatcher.',
    path: '/nodes',
    noindex: true,
  })

  useEffect(() => {
    fetchNodes()
  }, [])

  const fetchNodes = async () => {
    try {
      const data = await monitorService.getNodes()
      setNodes(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch nodes')
    } finally {
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'ready':
        return 'success'
      case 'notready':
        return 'error'
      default:
        return 'warning'
    }
  }

  const columns: GridColDef[] = [
    {
      field: 'name',
      headerName: 'Name',
      width: 200,
      renderHeader: () => <strong>Name</strong>,
    },
    {
      field: 'status',
      headerName: 'Status',
      width: 120,
      renderCell: (params) => (
        <Chip
          label={params.value}
          color={getStatusColor(params.value)}
          size="small"
        />
      ),
      renderHeader: () => <strong>Status</strong>,
    },
    {
      field: 'os',
      headerName: 'OS',
      width: 100,
      renderHeader: () => <strong>OS</strong>,
    },
    {
      field: 'architecture',
      headerName: 'Architecture',
      width: 120,
      renderHeader: () => <strong>Architecture</strong>,
    },
    {
      field: 'cpu',
      headerName: 'CPU',
      width: 100,
      renderHeader: () => <strong>CPU</strong>,
    },
    {
      field: 'memory',
      headerName: 'Memory',
      width: 100,
      renderHeader: () => <strong>Memory</strong>,
    },
    {
      field: 'pods',
      headerName: 'Pods',
      width: 100,
      renderHeader: () => <strong>Pods</strong>,
    },
    {
      field: 'kubelet_version',
      headerName: 'Kubelet Version',
      width: 150,
      renderHeader: () => <strong>Kubelet Version</strong>,
    },
    {
      field: 'created_at',
      headerName: 'Created',
      width: 180,
      renderCell: (params) => new Date(params.value).toLocaleString(),
      renderHeader: () => <strong>Created</strong>,
    },
  ]

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

  return (
    <Box>
      <Typography variant="h4" gutterBottom sx={{ color: 'white', mb: 3 }}>
        Nodes
      </Typography>

      <Box
        sx={{
          height: 600,
          width: '100%',
          '& .MuiDataGrid-root': {
            backgroundColor: '#1e1e1e',
            color: 'white',
            border: '1px solid #333',
          },
          '& .MuiDataGrid-cell': {
            borderBottom: '1px solid #333',
          },
          '& .MuiDataGrid-columnHeaders': {
            backgroundColor: '#2a2a2a',
            borderBottom: '2px solid #2196f3',
          },
          '& .MuiDataGrid-columnHeaderTitle': {
            color: 'white',
            fontWeight: 'bold',
          },
          '& .MuiDataGrid-row:hover': {
            backgroundColor: '#2a2a2a',
          },
          '& .MuiDataGrid-footerContainer': {
            backgroundColor: '#2a2a2a',
            color: 'white',
          },
          '& .MuiTablePagination-root': {
            color: 'white',
          },
          '& .MuiIconButton-root': {
            color: 'white',
          },
        }}
      >
        <DataGrid
          rows={nodes}
          columns={columns}
          getRowId={(row) => row.name}
          initialState={{
            pagination: {
              paginationModel: { pageSize: 10 },
            },
          }}
          pageSizeOptions={[10, 25, 50]}
          disableRowSelectionOnClick
        />
      </Box>
    </Box>
  )
}

export default Nodes