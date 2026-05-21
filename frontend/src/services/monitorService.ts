import api from './api'

export interface Cluster {
  nodes: NodeSummary
  pods: PodSummary
  deployments: DeploymentSummary
  services: ServiceSummary
  timestamp: number
}

export interface NodeSummary {
  total: number
  ready: number
  not_ready: number
}

export interface PodSummary {
  total: number
  running: number
  pending: number
  failed: number
}

export interface DeploymentSummary {
  total: number
  ready: number
  not_ready: number
}

export interface ServiceSummary {
  total: number
  cluster_ip: number
  load_balancer: number
  node_port: number
}

export interface Node {
  name: string
  status: string
  os: string
  kubelet_version: string
  architecture: string
  cpu: string
  memory: string
  pods: string
  created_at: string
}

export interface Pod {
  name: string
  namespace: string
  status: string
  node: string
  restarts: number
  ip: string
  created_at: string
  containers: ContainerInfo[]
}

export interface ContainerInfo {
  name: string
  image: string
  status: string
}

export const monitorService = {
  getClusterOverview: async (): Promise<Cluster> => {
    const response = await api.get('/monitor/cluster')
    return response.data
  },

  getNodes: async (): Promise<Node[]> => {
    const response = await api.get('/monitor/nodes')
    return response.data
  },

  getPods: async (namespace?: string): Promise<Pod[]> => {
    const response = await api.get('/monitor/pods', {
      params: { namespace },
    })
    return response.data
  },
}