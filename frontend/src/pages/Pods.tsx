import React, { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  CircularProgress,
  Alert,
  Chip,
  TextField,
  InputAdornment,
} from '@mui/material'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import SearchIcon from '@mui/icons-material/Search'
import { monitorService, Pod } from '../services/monitorService'
import { useSEO } from '../hooks/useSEO'

const Pods: React.FC = () => {
  const [pods, setPods] = useState<Pod[]>([])
  const [filteredPods, setFilteredPods] = useState<Pod[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')

  useSEO({
    title: 'Cluster Pods',
    description: 'Inspect Kubernetes pod health, namespaces, and lifecycle status with KubeWatcher.',
    path: '/pods',
    noindex: true,
  })

  useEffect(() => {
    fetchPods()
  }, [])

  useEffect(() => {
    const filtered = pods.filter(
      (pod) =>
        pod.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        pod.namespace.toLowerCase().includes(searchTerm.toLowerCase()) ||
        pod.status.toLowerCase().includes(searchTerm.toLowerCase())
    )
    setFilteredPods(filtered)
  }, [pods, searchTerm])

  const fetchPods = async () => {
    try {
      const data = await monitorService.getPods()
      setPods(data)
      setFilteredPods(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch pods')
    } finally {
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'running':
        return 'success'
      case 'pending':
        return 'warning'
      case 'failed':
      case 'crashloopbackoff':
        return 'error'
      default:
        return 'default'
    }
  }

  const columns: GridColDef[] = [
    {
      field: 'name',
      headerName: 'Name',
      width: 250,
      renderHeader: () => <strong>Name</strong>,
    },
    {
      field: 'namespace',
      headerName: 'Namespace',
      width: 120,
      renderHeader: () => <strong>Namespace</strong>,
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
      field: 'node',
      headerName: 'Node',
      width: 150,
      renderHeader: () => <strong>Node</strong>,
    },
    {
      field: 'restarts',
      headerName: 'Restarts',
      width: 100,
      renderHeader: () => <strong>Restarts</strong>,
      renderCell: (params) => (
        <Typography
          sx={{
            color: params.value > 0 ? '#ff9800' : 'inherit',
          }}
        >
          {params.value}
        </Typography>
      ),
    },
    {
      field: 'ip',
      headerName: 'IP',
      width: 130,
      renderHeader: () => <strong>IP</strong>,
    },
    {
      field: 'containers',
      headerName: 'Containers',
      width: 120,
      renderCell: (params) => params.value?.length || 0,
      renderHeader: () => <strong>Containers</strong>,
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
        Pods
      </Typography>

      <Box mb={2}>
        <TextField
          placeholder="Search pods..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon sx={{ color: 'white' }} />
              </InputAdornment>
            ),
          }}
          sx={{
            width: 300,
            '& .MuiInputBase-input': { color: 'white' },
            '& .MuiOutlinedInput-root': {
              '& fieldset': { borderColor: 'white' },
              '&:hover fieldset': { borderColor: '#2196f3' },
              '&.Mui-focused fieldset': { borderColor: '#2196f3' },
            },
          }}
        />
      </Box>

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
          rows={filteredPods}
          columns={columns}
          getRowId={(row) => `${row.namespace}-${row.name}`}
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

export default Pods