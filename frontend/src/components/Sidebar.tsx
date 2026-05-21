import React from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Box,
  Button,
  Divider,
} from '@mui/material'
import {
  Dashboard as DashboardIcon,
  Storage as NodesIcon,
  ViewList as PodsIcon,
  People as EmployeesIcon,
  Assessment as ReviewsIcon,
  Logout as LogoutIcon,
  AdminPanelSettings as AdminIcon,
  Public as ExternalIcon,
} from '@mui/icons-material'
import { useAuth } from '../hooks/useAuth'

const drawerWidth = 240

const Sidebar: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { logout, isAdmin, user } = useAuth()

  const commonMenuItems = [
    { text: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
    { text: 'Nodes', icon: <NodesIcon />, path: '/nodes' },
    { text: 'Pods', icon: <PodsIcon />, path: '/pods' },
    { text: 'Employees', icon: <EmployeesIcon />, path: '/employees' },
    { text: 'Reviews', icon: <ReviewsIcon />, path: '/reviews' },
    { text: 'External Posts', icon: <ExternalIcon />, path: '/external' },
  ]

  const adminMenuItems = [
    { text: 'User Management', icon: <AdminIcon />, path: '/admin/users' },
  ]

  const menuItems = isAdmin() ? [...commonMenuItems, ...adminMenuItems] : commonMenuItems

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
          backgroundColor: '#1e1e1e',
          color: 'white',
        },
      }}
    >
      <Toolbar>
        <Typography variant="h6" noWrap component="div">
          KubeWatcher
        </Typography>
      </Toolbar>
      <Box sx={{ p: 2 }}>
        <Typography variant="body2" sx={{ color: 'rgba(255, 255, 255, 0.7)' }}>
          {user?.name}
        </Typography>
        <Typography variant="caption" sx={{ color: 'rgba(255, 255, 255, 0.5)' }}>
          Role: <span style={{ textTransform: 'uppercase' }}>{user?.role}</span>
        </Typography>
      </Box>
      <Divider sx={{ backgroundColor: 'rgba(255, 255, 255, 0.1)' }} />
      <Box sx={{ overflow: 'auto', flexGrow: 1 }}>
        <List>
          {menuItems.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => navigate(item.path)}
                sx={{
                  '&.Mui-selected': {
                    backgroundColor: 'rgba(33, 150, 243, 0.2)',
                    '&:hover': {
                      backgroundColor: 'rgba(33, 150, 243, 0.3)',
                    },
                  },
                }}
              >
                <ListItemIcon sx={{ color: 'white' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText primary={item.text} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Box>
      <Box sx={{ p: 2 }}>
        <Button
          fullWidth
          variant="outlined"
          startIcon={<LogoutIcon />}
          onClick={handleLogout}
          sx={{
            color: 'white',
            borderColor: 'white',
            '&:hover': {
              borderColor: 'white',
              backgroundColor: 'rgba(255, 255, 255, 0.1)',
            },
          }}
        >
          Logout
        </Button>
      </Box>
    </Drawer>
  )
}

export default Sidebar