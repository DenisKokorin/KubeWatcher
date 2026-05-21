import { useEffect, useState } from 'react'
import {
  Box,
  Button,
  CircularProgress,
  Grid,
  Paper,
  TextField,
  Typography,
} from '@mui/material'
import RefreshIcon from '@mui/icons-material/Refresh'
import SearchIcon from '@mui/icons-material/Search'
import { externalService, ExternalPost } from '../services/externalService'

const ExternalData = () => {
  const [posts, setPosts] = useState<ExternalPost[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchPosts = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await externalService.getPosts({ limit: 20, search })
      setPosts(data)
    } catch (err) {
      setError('Unable to load external posts. Please try again later.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchPosts()
  }, [])

  return (
    <Box>
      <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
        <Box>
          <Typography variant="h4" gutterBottom>
            External posts
          </Typography>
          <Typography color="text.secondary">
            Loaded from an external API adapter with resilient backend handling.
          </Typography>
        </Box>
        <Button variant="contained" startIcon={<RefreshIcon />} onClick={fetchPosts}>
          Reload
        </Button>
      </Box>

      <Box display="flex" gap={2} mb={3} flexWrap="wrap">
        <TextField
          label="Search posts"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(event) => {
            if (event.key === 'Enter') {
              fetchPosts()
            }
          }}
          size="small"
          InputProps={{
            endAdornment: <SearchIcon color="action" />,
          }}
        />
        <Button variant="outlined" onClick={fetchPosts}>
          Apply
        </Button>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="240px">
          <CircularProgress />
        </Box>
      ) : error ? (
        <Paper sx={{ p: 3, backgroundColor: '#fdecea', color: '#611a15' }}>
          <Typography variant="h6">External API unavailable</Typography>
          <Typography>{error}</Typography>
        </Paper>
      ) : posts.length === 0 ? (
        <Paper sx={{ p: 3 }}>
          <Typography>No external posts found.</Typography>
        </Paper>
      ) : (
        <Grid container spacing={2}>
          {posts.map((post) => (
            <Grid item xs={12} md={6} key={post.id}>
              <Paper sx={{ p: 3, height: '100%' }}>
                <Typography variant="h6" gutterBottom>
                  {post.title}
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  {post.summary}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Source: {post.source} • User ID: {post.user_id}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      )}
    </Box>
  )
}

export default ExternalData
