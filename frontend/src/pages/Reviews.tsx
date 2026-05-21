import React, { useState } from 'react'
import {
  Box,
  Typography,
  CircularProgress,
  Alert,
  Chip,
  Rating,
} from '@mui/material'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import { reviewService, TeamStats } from '../services/reviewService'
import { useSEO } from '../hooks/useSEO'

interface TeamData extends TeamStats {
  id: string
}

const Reviews: React.FC = () => {
  const [teams, setTeams] = useState<TeamData[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [teamInput, setTeamInput] = useState('')

  useSEO({
    title: 'Performance Reviews',
    description: 'Track team performance and review results inside the KubeWatcher performance dashboard.',
    path: '/reviews',
    noindex: true,
  })

  const fetchTeamStats = async (teamName: string) => {
    if (!teamName.trim()) {
      setError('Please enter a team name')
      return
    }
    
    setLoading(true)
    setError('')
    try {
      const data = await reviewService.getTeamStats(teamName)
      const teamData: TeamData = { ...data, id: teamName }
      
      // Add or update team in the list
      setTeams((prev) => {
        const existing = prev.findIndex((t) => t.id === teamName)
        if (existing >= 0) {
          const updated = [...prev]
          updated[existing] = teamData
          return updated
        }
        return [...prev, teamData]
      })
      setTeamInput('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch team stats')
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    fetchTeamStats(teamInput)
  }

  const columns: GridColDef[] = [
    {
      field: 'id',
      headerName: 'Team',
      width: 150,
      renderHeader: () => <strong>Team</strong>,
    },
    {
      field: 'employee_count',
      headerName: 'Employees',
      width: 120,
      renderHeader: () => <strong>Employees</strong>,
    },
    {
      field: 'reviews_count',
      headerName: 'Total Reviews',
      width: 150,
      renderHeader: () => <strong>Total Reviews</strong>,
    },
    {
      field: 'average_rating',
      headerName: 'Average Rating',
      width: 180,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Rating
            value={Math.round(params.value * 10) / 10}
            readOnly
            size="small"
            precision={0.1}
          />
          <Typography variant="body2">{params.value.toFixed(2)}</Typography>
        </Box>
      ),
      renderHeader: () => <strong>Average Rating</strong>,
    },
    {
      field: 'rating',
      headerName: 'Rating',
      width: 120,
      renderCell: (params) => (
        <Rating value={params.value} readOnly size="small" />
      ),
      renderHeader: () => <strong>Rating</strong>,
    },
    {
      field: 'comments',
      headerName: 'Comments',
      width: 300,
      renderCell: (params) => (
        <Typography
          sx={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {params.value}
        </Typography>
      ),
      renderHeader: () => <strong>Comments</strong>,
    },
    {
      field: 'goals',
      headerName: 'Goals',
      width: 200,
      renderCell: (params) => (
        <Box>
          {params.value?.slice(0, 2).map((goal: string, index: number) => (
            <Chip
              key={index}
              label={goal}
              size="small"
              sx={{
                mr: 0.5,
                mb: 0.5,
                backgroundColor: '#2196f3',
                color: 'white',
              }}
            />
          ))}
          {params.value?.length > 2 && (
            <Typography variant="caption" sx={{ color: 'gray' }}>
              +{params.value.length - 2} more
            </Typography>
          )}
        </Box>
      ),
      renderHeader: () => <strong>Goals</strong>,
    },
    {
      field: 'strengths',
      headerName: 'Strengths',
      width: 200,
      renderCell: (params) => (
        <Box>
          {params.value?.slice(0, 2).map((strength: string, index: number) => (
            <Chip
              key={index}
              label={strength}
              size="small"
              sx={{
                mr: 0.5,
                mb: 0.5,
                backgroundColor: '#4caf50',
                color: 'white',
              }}
            />
          ))}
          {params.value?.length > 2 && (
            <Typography variant="caption" sx={{ color: 'gray' }}>
              +{params.value.length - 2} more
            </Typography>
          )}
        </Box>
      ),
      renderHeader: () => <strong>Strengths</strong>,
    },
    {
      field: 'areas_for_improvement',
      headerName: 'Areas for Improvement',
      width: 200,
      renderCell: (params) => (
        <Box>
          {params.value?.slice(0, 2).map((area: string, index: number) => (
            <Chip
              key={index}
              label={area}
              size="small"
              sx={{
                mr: 0.5,
                mb: 0.5,
                backgroundColor: '#ff9800',
                color: 'white',
              }}
            />
          ))}
          {params.value?.length > 2 && (
            <Typography variant="caption" sx={{ color: 'gray' }}>
              +{params.value.length - 2} more
            </Typography>
          )}
        </Box>
      ),
      renderHeader: () => <strong>Areas for Improvement</strong>,
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
        Reviews
      </Typography>

      <Box sx={{ mb: 3, display: 'flex', gap: 1 }}>
        <input
          type="text"
          placeholder="Search team..."
          value={teamInput}
          onChange={(e) => setTeamInput(e.target.value)}
          style={{
            flex: 1,
            padding: '12px',
            backgroundColor: '#2a2a2a',
            color: 'white',
            border: '1px solid #444',
            borderRadius: '4px',
            fontSize: '14px',
          }}
        />
        <button
          onClick={handleSearch}
          style={{
            padding: '12px 24px',
            backgroundColor: '#2196f3',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '14px',
            fontWeight: 'bold',
          }}
        >
          Search
        </button>
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
          rows={teams}
          columns={columns}
          getRowId={(row) => row.id}
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

export default Reviews