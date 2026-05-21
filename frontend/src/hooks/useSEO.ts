import { useEffect } from 'react'

interface SEOOptions {
  title: string
  description: string
  path?: string
  noindex?: boolean
}

const setTag = (selector: string, value: string) => {
  let element = document.head.querySelector(selector) as HTMLElement | null
  if (!element) {
    element = document.createElement('meta')
    const [attrName, attrValue] = selector.split('[')[1].replace(']', '').split('=')
    element.setAttribute(attrName, attrValue.replace(/"/g, ''))
    document.head.appendChild(element)
  }
  element.setAttribute('content', value)
}

const setLink = (rel: string, href: string) => {
  let link = document.head.querySelector(`link[rel="${rel}"]`) as HTMLLinkElement | null
  if (!link) {
    link = document.createElement('link')
    link.rel = rel
    document.head.appendChild(link)
  }
  link.href = href
}

export const useSEO = ({ title, description, path = '/', noindex = false }: SEOOptions) => {
  useEffect(() => {
    const formattedTitle = `${title} | KubeWatcher`
    document.title = formattedTitle

    setTag('meta[name="description"]', description)
    setTag('meta[name="robots"]', noindex ? 'noindex, nofollow' : 'index, follow')
    setTag('meta[property="og:title"]', formattedTitle)
    setTag('meta[property="og:description"]', description)
    setTag('meta[property="og:type"]', 'website')
    setTag('meta[property="og:site_name"]', 'KubeWatcher')
    setTag('meta[property="og:url"]', `${window.location.origin}${path}`)
    setLink('canonical', `${window.location.origin}${path}`)
  }, [title, description, path, noindex])
}
