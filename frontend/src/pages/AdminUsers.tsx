import React, { useState, useEffect } from 'react'
import {
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Typography,
  Alert,
  CircularProgress,
  Chip,
  Box,
} from '@mui/material'
import api from '../services/api'
import { useSEO } from '../hooks/useSEO'

interface User {
  uuid: string
  name: string
  email: string
  position: string
  team: string
  role: string
  hire_date: string
  profile_document_url?: string
}

const AdminUsers: React.FC = () => {
  useSEO({
    title: 'Admin User Management',
    description: 'Manage users, roles and secure document workflows in the KubeWatcher admin panel.',
    path: '/admin/users',
    noindex: true,
  })
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [newRole, setNewRole] = useState<string>('')
  const [openDialog, setOpenDialog] = useState(false)
  const [updating, setUpdating] = useState(false)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [newUser, setNewUser] = useState({ name: '', email: '', position: '', team: '', role: 'user', password: '' })
  const [creating, setCreating] = useState(false)
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false)
  const [selectedUserForUpload, setSelectedUserForUpload] = useState<User | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)
  const [downloadDialogOpen, setDownloadDialogOpen] = useState(false)
  const [selectedUserForDownload, setSelectedUserForDownload] = useState<User | null>(null)
  const [downloading, setDownloading] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedUserForDelete, setSelectedUserForDelete] = useState<User | null>(null)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const response = await api.get('/admin/users')
      setUsers(response.data.users || [])
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch users')
    } finally {
      setLoading(false)
    }
  }

  const handleOpenDialog = (user: User) => {
    setSelectedUser(user)
    setNewRole(user.role)
    setOpenDialog(true)
  }

  const handleCloseDialog = () => {
    setOpenDialog(false)
    setSelectedUser(null)
    setNewRole('')
  }

  const handleUpdateRole = async () => {
    if (!selectedUser || !newRole) return

    try {
      setUpdating(true)
      await api.put(`/admin/users/${selectedUser.uuid}/role`, { role: newRole })
      setUsers((prevUsers) =>
        prevUsers.map((u) =>
          u.uuid === selectedUser.uuid ? { ...u, role: newRole } : u
        )
      )
      handleCloseDialog()
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update role')
    } finally {
      setUpdating(false)
    }
  }

  const handleOpenCreateDialog = () => {
    setCreateDialogOpen(true)
    setNewUser({ name: '', email: '', position: '', team: '', role: 'user', password: '' })
  }

  const handleCloseCreateDialog = () => {
    setCreateDialogOpen(false)
  }

  const handleCreateUser = async () => {
    if (!newUser.name || !newUser.email || !newUser.password) {
      setError('Name, email, and password are required')
      return
    }

    try {
      setCreating(true)
      await api.post('/admin/users', newUser)
      await fetchUsers()
      handleCloseCreateDialog()
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create user')
    } finally {
      setCreating(false)
    }
  }

  const handleOpenUploadDialog = (user: User) => {
    setSelectedUserForUpload(user)
    setSelectedFile(null)
    setUploadDialogOpen(true)
  }

  const handleCloseUploadDialog = () => {
    setUploadDialogOpen(false)
    setSelectedUserForUpload(null)
    setSelectedFile(null)
  }

  const handleUploadDocument = async () => {
    if (!selectedUserForUpload || !selectedFile) {
      setError('Please select a user and a file before uploading')
      return
    }

    // Validate file
    const validationError = validateFile(selectedFile)
    if (validationError) {
      setError(validationError)
      return
    }

    try {
      setUploading(true)
      const formData = new FormData()
      formData.append('file', selectedFile)
      const response = await api.post(`/admin/users/${selectedUserForUpload.uuid}/document`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      setError('')
      handleCloseUploadDialog()
      await fetchUsers()
      console.log('Document uploaded:', response.data)
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to upload document'
      setError(errorMessage)
    } finally {
      setUploading(false)
    }
  }

  const validateFile = (file: File): string | null => {
    const maxSize = 10 * 1024 * 1024 // 10MB
    const allowedTypes = ['.pdf', '.doc', '.docx', '.txt', '.jpg', '.jpeg', '.png', '.gif']

    if (file.size > maxSize) {
      return 'File size must be less than 10MB'
    }

    const fileExt = file.name.toLowerCase().substring(file.name.lastIndexOf('.'))
    if (!allowedTypes.includes(fileExt)) {
      return 'File type not allowed. Allowed types: PDF, DOC, DOCX, TXT, JPG, PNG, GIF'
    }

    return null
  }

  const handleOpenDownloadDialog = (user: User) => {
    setSelectedUserForDownload(user)
    setDownloadDialogOpen(true)
  }

  const handleCloseDownloadDialog = () => {
    setDownloadDialogOpen(false)
    setSelectedUserForDownload(null)
  }

  const handleDownloadDocument = async () => {
    if (!selectedUserForDownload) return

    try {
      setDownloading(true)
      const response = await api.get(`/admin/users/${selectedUserForDownload.uuid}/document`)
      const presignedUrl = response.data.data.url

      // Open the presigned URL in a new tab
      window.open(presignedUrl, '_blank')
      handleCloseDownloadDialog()
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to get document URL')
    } finally {
      setDownloading(false)
    }
  }

  const handleOpenDeleteDialog = (user: User) => {
    setSelectedUserForDelete(user)
    setDeleteDialogOpen(true)
  }

  const handleCloseDeleteDialog = () => {
    setDeleteDialogOpen(false)
    setSelectedUserForDelete(null)
  }

  const handleDeleteDocument = async () => {
    if (!selectedUserForDelete) return

    try {
      setDeleting(true)
      await api.delete(`/admin/users/${selectedUserForDelete.uuid}/document`)

      // Update the user in the list to remove the document URL
      setUsers((prevUsers) =>
        prevUsers.map((u) =>
          u.uuid === selectedUserForDelete.uuid
            ? { ...u, profile_document_url: undefined }
            : u
        )
      )

      handleCloseDeleteDialog()
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete document')
    } finally {
      setDeleting(false)
    }
  }

  const getRoleColor = (role: string) => {
    return role === 'admin' ? 'error' : 'default'
  }

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center">
          <CircularProgress />
        </Box>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h4" component="h1" sx={{ color: 'white' }}>
          User Management
        </Typography>
        <Button variant="contained" color="primary" onClick={handleOpenCreateDialog}>
          Create User
        </Button>
      </Box>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
              <TableCell sx={{ fontWeight: 'bold' }}>Name</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Email</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Position</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Team</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Role</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Document</TableCell>
              <TableCell sx={{ fontWeight: 'bold' }}>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.uuid} hover>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>{user.position}</TableCell>
                <TableCell>{user.team}</TableCell>
                <TableCell>
                  <Chip
                    label={user.role.toUpperCase()}
                    color={getRoleColor(user.role)}
                    variant="outlined"
                  />
                </TableCell>
                <TableCell>
                  {user.profile_document_url ? (
                    <Chip
                      label="Has Document"
                      color="success"
                      variant="outlined"
                      size="small"
                    />
                  ) : (
                    <Chip
                      label="No Document"
                      color="default"
                      variant="outlined"
                      size="small"
                    />
                  )}
                </TableCell>
                <TableCell sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => handleOpenDialog(user)}
                  >
                    Change Role
                  </Button>
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => handleOpenUploadDialog(user)}
                  >
                    Upload Document
                  </Button>
                  {user.profile_document_url && (
                    <>
                      <Button
                        variant="outlined"
                        size="small"
                        color="primary"
                        onClick={() => handleOpenDownloadDialog(user)}
                      >
                        Download
                      </Button>
                      <Button
                        variant="outlined"
                        size="small"
                        color="error"
                        onClick={() => handleOpenDeleteDialog(user)}
                      >
                        Delete Doc
                      </Button>
                    </>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Change User Role</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {selectedUser && (
            <>
              <p>
                <strong>User:</strong> {selectedUser.name} ({selectedUser.email})
              </p>
              <FormControl fullWidth sx={{ mt: 2 }}>
                <InputLabel>Role</InputLabel>
                <Select
                  value={newRole}
                  label="Role"
                  onChange={(e) => setNewRole(e.target.value)}
                >
                  <MenuItem value="user">User</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </Select>
              </FormControl>
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button
            onClick={handleUpdateRole}
            variant="contained"
            disabled={updating || !newRole || newRole === selectedUser?.role}
          >
            {updating ? <CircularProgress size={24} /> : 'Update'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={createDialogOpen} onClose={handleCloseCreateDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Create User</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            margin="dense"
            label="Name"
            fullWidth
            value={newUser.name}
            onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Email"
            fullWidth
            value={newUser.email}
            onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Position"
            fullWidth
            value={newUser.position}
            onChange={(e) => setNewUser({ ...newUser, position: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Team"
            fullWidth
            value={newUser.team}
            onChange={(e) => setNewUser({ ...newUser, team: e.target.value })}
          />
          <FormControl fullWidth sx={{ mt: 2 }}>
            <InputLabel>Role</InputLabel>
            <Select
              value={newUser.role}
              label="Role"
              onChange={(e) => setNewUser({ ...newUser, role: e.target.value })}
            >
              <MenuItem value="user">User</MenuItem>
              <MenuItem value="admin">Admin</MenuItem>
            </Select>
          </FormControl>
          <TextField
            margin="dense"
            label="Password"
            type="password"
            fullWidth
            value={newUser.password}
            onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCreateDialog}>Cancel</Button>
          <Button
            onClick={handleCreateUser}
            variant="contained"
            disabled={creating || !newUser.name || !newUser.email || !newUser.password}
          >
            {creating ? <CircularProgress size={24} /> : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={uploadDialogOpen} onClose={handleCloseUploadDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Upload Document</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {selectedUserForUpload && (
            <Box sx={{ mb: 2 }}>
              <Typography><strong>User:</strong> {selectedUserForUpload.name}</Typography>
              <Typography><strong>Email:</strong> {selectedUserForUpload.email}</Typography>
            </Box>
          )}
          <Button variant="contained" component="label">
            Choose File
            <input
              type="file"
              hidden
              onChange={(e) => {
                const file = e.target.files?.[0]
                if (file) {
                  setSelectedFile(file)
                }
              }}
            />
          </Button>
          {selectedFile && (
            <Typography sx={{ mt: 2 }}>{selectedFile.name}</Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUploadDialog}>Cancel</Button>
          <Button
            onClick={handleUploadDocument}
            variant="contained"
            disabled={uploading || !selectedFile}
          >
            {uploading ? <CircularProgress size={24} /> : 'Upload'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={downloadDialogOpen} onClose={handleCloseDownloadDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Download Document</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {selectedUserForDownload && (
            <Box sx={{ mb: 2 }}>
              <Typography><strong>User:</strong> {selectedUserForDownload.name}</Typography>
              <Typography><strong>Email:</strong> {selectedUserForDownload.email}</Typography>
              <Typography sx={{ mt: 2 }}>
                Clicking "Download" will open the document in a new tab.
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDownloadDialog}>Cancel</Button>
          <Button
            onClick={handleDownloadDocument}
            variant="contained"
            disabled={downloading}
          >
            {downloading ? <CircularProgress size={24} /> : 'Download'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={deleteDialogOpen} onClose={handleCloseDeleteDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Delete Document</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {selectedUserForDelete && (
            <Box sx={{ mb: 2 }}>
              <Typography><strong>User:</strong> {selectedUserForDelete.name}</Typography>
              <Typography><strong>Email:</strong> {selectedUserForDelete.email}</Typography>
              <Typography sx={{ mt: 2, color: 'error.main' }}>
                Are you sure you want to delete this user's document? This action cannot be undone.
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog}>Cancel</Button>
          <Button
            onClick={handleDeleteDocument}
            variant="contained"
            color="error"
            disabled={deleting}
          >
            {deleting ? <CircularProgress size={24} /> : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}

export default AdminUsers
