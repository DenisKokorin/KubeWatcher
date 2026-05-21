import React, { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  CircularProgress,
  Alert,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Rating,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material'
import { DataGrid, GridColDef, GridSortModel } from '@mui/x-data-grid'
import { reviewService, Employee, EmployeeWithReviews, ReviewRequest } from '../services/reviewService'
import { useAuth } from '../hooks/useAuth'
import { useSEO } from '../hooks/useSEO'

const Employees: React.FC = () => {
  const { user } = useAuth()
  const [employees, setEmployees] = useState<Employee[]>([])
  const [totalRows, setTotalRows] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [openDialog, setOpenDialog] = useState(false)
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(null)
  const [viewEmployee, setViewEmployee] = useState<EmployeeWithReviews | null>(null)
  const [openViewDialog, setOpenViewDialog] = useState(false)
  const [showProfile, setShowProfile] = useState(false)
  const [search, setSearch] = useState('')
  const [teamFilter, setTeamFilter] = useState('')
  const [positionFilter, setPositionFilter] = useState('')
  const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: 10 })
  const [sortingModel, setSortingModel] = useState<GridSortModel>([{ field: 'name', sort: 'asc' }])
  const [reviewData, setReviewData] = useState<Partial<ReviewRequest>>({
    rating: 0,
    comments: '',
    goals: [],
    strengths: [],
    areas_for_improvement: [],
  })

  useSEO({
    title: 'Employee Review Management',
    description: 'Manage employee profiles, reviews, and team insights securely in KubeWatcher.',
    path: '/employees',
    noindex: true,
  })

  useSEO({
    title: 'Employee Review Management',
    description: 'Access employee profiles and review workflows within the KubeWatcher application.',
    path: '/employees',
    noindex: true,
  })

  useEffect(() => {
    fetchEmployees()
  }, [search, teamFilter, positionFilter, paginationModel.page, paginationModel.pageSize, sortingModel])

  const fetchEmployees = async () => {
    setLoading(true)
    try {
      const sortField = sortingModel[0]?.field || 'name'
      const sortOrder = sortingModel[0]?.sort || 'asc'
      const data = await reviewService.getAllEmployees({
        search,
        team: teamFilter,
        position: positionFilter,
        page: paginationModel.page,
        size: paginationModel.pageSize,
        sortBy: sortField,
        sortOrder,
      })
      setEmployees(data.employees)
      setTotalRows(data.total)
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch employees')
    } finally {
      setLoading(false)
    }
  }

  const handleViewEmployee = async (employee: Employee) => {
    try {
      const result = await reviewService.getEmployeeWithReviews(employee.uuid)
      setViewEmployee(result)
      setOpenViewDialog(true)
      setShowProfile(false)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load employee details')
    }
  }

  const handleShowMyProfile = () => {
    setOpenViewDialog(true)
    setShowProfile(true)
    setViewEmployee(null)
  }

  const handleCreateReview = async () => {
    if (!selectedEmployee || !reviewData.rating || !user) return

    console.log('Selected employee:', selectedEmployee)
    console.log('Current user:', user)
    console.log('Review data:', reviewData)

    try {
      await reviewService.createReview({
        employee_id: selectedEmployee.uuid,
        reviewer_id: user.uuid,
        period: new Date().getFullYear() + ' Q' + Math.ceil((new Date().getMonth() + 1) / 3),
        rating: reviewData.rating,
        comments: reviewData.comments || '',
        goals: reviewData.goals || [],
        strengths: reviewData.strengths || [],
        areas_for_improvement: reviewData.areas_for_improvement || [],
      })
      setOpenDialog(false)
      setReviewData({
        rating: 0,
        comments: '',
        goals: [],
        strengths: [],
        areas_for_improvement: [],
      })
      // Refresh data if needed
      fetchEmployees()
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create review')
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
      field: 'email',
      headerName: 'Email',
      width: 250,
      renderHeader: () => <strong>Email</strong>,
    },
    {
      field: 'position',
      headerName: 'Position',
      width: 150,
      renderHeader: () => <strong>Position</strong>,
    },
    {
      field: 'team',
      headerName: 'Team',
      width: 120,
      renderHeader: () => <strong>Team</strong>,
    },
    {
      field: 'role',
      headerName: 'Role',
      width: 120,
      renderHeader: () => <strong>Role</strong>,
    },
    {
      field: 'hire_date',
      headerName: 'Hire Date',
      width: 120,
      renderCell: (params) => new Date(params.value).toLocaleDateString(),
      renderHeader: () => <strong>Hire Date</strong>,
    },
    {
      field: 'profile_document_url',
      headerName: 'Document',
      width: 120,
      renderCell: (params) => (
        params.value ? (
          <Typography sx={{ color: 'success.main', fontWeight: 'bold' }}>
            ✓ Available
          </Typography>
        ) : (
          <Typography sx={{ color: 'text.secondary' }}>
            No Document
          </Typography>
        )
      ),
      renderHeader: () => <strong>Document</strong>,
    },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 220,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            size="small"
            onClick={() => handleViewEmployee(params.row)}
            sx={{
              color: '#4caf50',
              borderColor: '#4caf50',
              '&:hover': { borderColor: '#4caf50', backgroundColor: 'rgba(76, 175, 80, 0.1)' },
            }}
          >
            View
          </Button>

          <Button
            variant="outlined"
            size="small"
            onClick={() => {
              setSelectedEmployee(params.row)
              setOpenDialog(true)
            }}
            sx={{
              color: '#2196f3',
              borderColor: '#2196f3',
              '&:hover': { borderColor: '#2196f3', backgroundColor: 'rgba(33, 150, 243, 0.1)' },
            }}
          >
            Add Review
          </Button>
        </Box>
      ),
      renderHeader: () => <strong>Actions</strong>,
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
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4" gutterBottom sx={{ color: 'white' }}>
          Employees
        </Typography>
        {user && (
          <Button variant="outlined" onClick={handleShowMyProfile} sx={{ color: 'white', borderColor: 'white' }}>
            View My Profile
          </Button>
        )}
      </Box>

      <Box sx={{ mb: 2, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
        <TextField
          label="Search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          sx={{ minWidth: 240, backgroundColor: '#121212' }}
          InputLabelProps={{ style: { color: 'white' } }}
          InputProps={{ style: { color: 'white' } }}
        />
        <FormControl sx={{ minWidth: 180, backgroundColor: '#121212' }}>
          <InputLabel sx={{ color: 'white' }}>Team</InputLabel>
          <Select
            value={teamFilter}
            label="Team"
            onChange={(e) => setTeamFilter(e.target.value)}
            sx={{ color: 'white' }}
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="Engineering">Engineering</MenuItem>
            <MenuItem value="Product">Product</MenuItem>
            <MenuItem value="Design">Design</MenuItem>
          </Select>
        </FormControl>
        <FormControl sx={{ minWidth: 180, backgroundColor: '#121212' }}>
          <InputLabel sx={{ color: 'white' }}>Position</InputLabel>
          <Select
            value={positionFilter}
            label="Position"
            onChange={(e) => setPositionFilter(e.target.value)}
            sx={{ color: 'white' }}
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="Engineer">Engineer</MenuItem>
            <MenuItem value="Manager">Manager</MenuItem>
            <MenuItem value="Director">Director</MenuItem>
          </Select>
        </FormControl>
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
          rows={employees}
          columns={columns}
          getRowId={(row) => row.uuid}
          paginationMode="server"
          sortingMode="server"
          rowCount={totalRows}
          paginationModel={paginationModel}
          onPaginationModelChange={(model) => setPaginationModel(model)}
          sortModel={sortingModel}
          onSortModelChange={(model) => setSortingModel(model)}
          pageSizeOptions={[10, 25, 50]}
          disableRowSelectionOnClick
        />
      </Box>

      <Dialog
        open={openViewDialog}
        onClose={() => {
          setOpenViewDialog(false)
          setShowProfile(false)
          setViewEmployee(null)
        }}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: {
            backgroundColor: '#1e1e1e',
            color: 'white',
          },
        }}
      >
        <DialogTitle>
          {showProfile ? 'My Profile' : `Employee: ${viewEmployee?.employee.name || ''}`}
        </DialogTitle>
        <DialogContent>
          {showProfile && user && (
            <Box sx={{ color: 'white' }}>
              <Typography><strong>Name:</strong> {user.name}</Typography>
              <Typography><strong>Email:</strong> {user.email}</Typography>
              <Typography><strong>Position:</strong> {user.position}</Typography>
              <Typography><strong>Team:</strong> {user.team}</Typography>
            </Box>
          )}
          {!showProfile && viewEmployee && (
            <Box>
              <Typography><strong>Name:</strong> {viewEmployee.employee.name}</Typography>
              <Typography><strong>Email:</strong> {viewEmployee.employee.email}</Typography>
              <Typography><strong>Position:</strong> {viewEmployee.employee.position}</Typography>
              <Typography><strong>Team:</strong> {viewEmployee.employee.team}</Typography>
              <Typography><strong>Hire Date:</strong> {new Date(viewEmployee.employee.hire_date).toLocaleDateString()}</Typography>

              <Typography sx={{ mt: 2, mb: 1 }}><strong>Reviews</strong></Typography>
              {viewEmployee.reviews.length === 0 ? (
                <Typography>No reviews yet</Typography>
              ) : (
                viewEmployee.reviews.map((r) => (
                  <Box key={r.id} sx={{ mb: 1, p: 1, border: '1px solid #333', borderRadius: 1 }}>
                    <Typography><strong>Period:</strong> {r.period}</Typography>
                    <Typography><strong>Rating:</strong> {r.rating}</Typography>
                    <Typography><strong>Comments:</strong> {r.comments}</Typography>
                    <Typography><strong>Created:</strong> {new Date(r.created_at).toLocaleString()}</Typography>
                  </Box>
                ))
              )}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenViewDialog(false)} sx={{ color: 'white' }}>
            Close
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={openDialog}
        onClose={() => setOpenDialog(false)}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: {
            backgroundColor: '#1e1e1e',
            color: 'white',
          },
        }}
      >
        <DialogTitle>
          Add Review for {selectedEmployee?.name}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography component="legend">Rating</Typography>
            <Rating
              value={reviewData.rating}
              onChange={(_event, newValue) => {
                setReviewData({ ...reviewData, rating: newValue || 0 })
              }}
              sx={{
                '& .MuiRating-icon': {
                  color: '#2196f3',
                },
              }}
            />
          </Box>

          <TextField
            margin="dense"
            label="Comments"
            fullWidth
            multiline
            rows={3}
            value={reviewData.comments}
            onChange={(e) => setReviewData({ ...reviewData, comments: e.target.value })}
            sx={{
              mt: 2,
              '& .MuiInputBase-input': { color: 'white' },
              '& .MuiInputLabel-root': { color: 'white' },
              '& .MuiOutlinedInput-root': {
                '& fieldset': { borderColor: 'white' },
                '&:hover fieldset': { borderColor: '#2196f3' },
                '&.Mui-focused fieldset': { borderColor: '#2196f3' },
              },
            }}
          />

          <TextField
            margin="dense"
            label="Goals (comma-separated)"
            fullWidth
            value={reviewData.goals?.join(', ')}
            onChange={(e) => setReviewData({
              ...reviewData,
              goals: e.target.value.split(',').map(g => g.trim()).filter(g => g)
            })}
            sx={{
              mt: 2,
              '& .MuiInputBase-input': { color: 'white' },
              '& .MuiInputLabel-root': { color: 'white' },
              '& .MuiOutlinedInput-root': {
                '& fieldset': { borderColor: 'white' },
                '&:hover fieldset': { borderColor: '#2196f3' },
                '&.Mui-focused fieldset': { borderColor: '#2196f3' },
              },
            }}
          />

          <TextField
            margin="dense"
            label="Strengths (comma-separated)"
            fullWidth
            value={reviewData.strengths?.join(', ')}
            onChange={(e) => setReviewData({
              ...reviewData,
              strengths: e.target.value.split(',').map(s => s.trim()).filter(s => s)
            })}
            sx={{
              mt: 2,
              '& .MuiInputBase-input': { color: 'white' },
              '& .MuiInputLabel-root': { color: 'white' },
              '& .MuiOutlinedInput-root': {
                '& fieldset': { borderColor: 'white' },
                '&:hover fieldset': { borderColor: '#2196f3' },
                '&.Mui-focused fieldset': { borderColor: '#2196f3' },
              },
            }}
          />

          <TextField
            margin="dense"
            label="Areas for Improvement (comma-separated)"
            fullWidth
            value={reviewData.areas_for_improvement?.join(', ')}
            onChange={(e) => setReviewData({
              ...reviewData,
              areas_for_improvement: e.target.value.split(',').map(a => a.trim()).filter(a => a)
            })}
            sx={{
              mt: 2,
              '& .MuiInputBase-input': { color: 'white' },
              '& .MuiInputLabel-root': { color: 'white' },
              '& .MuiOutlinedInput-root': {
                '& fieldset': { borderColor: 'white' },
                '&:hover fieldset': { borderColor: '#2196f3' },
                '&.Mui-focused fieldset': { borderColor: '#2196f3' },
              },
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)} sx={{ color: 'white' }}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateReview}
            variant="contained"
            disabled={!reviewData.rating}
          >
            Create Review
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}

export default Employees