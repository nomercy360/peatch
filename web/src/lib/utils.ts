import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function detectLinkType(url: string): { type: string; icon: string } {
  const hostname = new URL(url).hostname.toLowerCase()

  // Common social media and professional platforms
  if (hostname.includes('github.com')) {
    return { type: 'github', icon: 'code' }
  }
  if (hostname.includes('linkedin.com')) {
    return { type: 'linkedin', icon: 'business_center' }
  }
  if (hostname.includes('twitter.com') || hostname.includes('x.com')) {
    return { type: 'twitter', icon: 'tag' }
  }
  if (hostname.includes('instagram.com')) {
    return { type: 'instagram', icon: 'photo_camera' }
  }
  if (hostname.includes('youtube.com') || hostname.includes('youtu.be')) {
    return { type: 'youtube', icon: 'play_circle' }
  }
  if (hostname.includes('behance.net')) {
    return { type: 'behance', icon: 'palette' }
  }
  if (hostname.includes('dribbble.com')) {
    return { type: 'dribbble', icon: 'brush' }
  }
  if (hostname.includes('medium.com')) {
    return { type: 'medium', icon: 'article' }
  }
  if (hostname.includes('dev.to')) {
    return { type: 'dev', icon: 'terminal' }
  }
  if (hostname.includes('stackoverflow.com')) {
    return { type: 'stackoverflow', icon: 'help' }
  }
  if (hostname.includes('telegram') || hostname.includes('t.me')) {
    return { type: 'telegram', icon: 'send' }
  }
  if (hostname.includes('discord')) {
    return { type: 'discord', icon: 'forum' }
  }
  if (hostname.includes('facebook.com') || hostname.includes('fb.com')) {
    return { type: 'facebook', icon: 'group' }
  }
  if (hostname.includes('tiktok.com')) {
    return { type: 'tiktok', icon: 'music_note' }
  }
  if (hostname.includes('pinterest.com')) {
    return { type: 'pinterest', icon: 'push_pin' }
  }
  if (hostname.includes('reddit.com')) {
    return { type: 'reddit', icon: 'forum' }
  }
  if (hostname.includes('twitch.tv')) {
    return { type: 'twitch', icon: 'videogame_asset' }
  }
  if (hostname.includes('spotify.com')) {
    return { type: 'spotify', icon: 'headphones' }
  }
  if (hostname.includes('soundcloud.com')) {
    return { type: 'soundcloud', icon: 'cloud' }
  }
  if (hostname.includes('notion.so')) {
    return { type: 'notion', icon: 'note' }
  }
  if (hostname.includes('figma.com')) {
    return { type: 'figma', icon: 'design_services' }
  }
  if (hostname.includes('gitlab.com')) {
    return { type: 'gitlab', icon: 'merge' }
  }
  if (hostname.includes('bitbucket.org')) {
    return { type: 'bitbucket', icon: 'commit' }
  }
  if (hostname.includes('codepen.io')) {
    return { type: 'codepen', icon: 'edit' }
  }
  if (hostname.includes('codesandbox.io')) {
    return { type: 'codesandbox', icon: 'web' }
  }
  if (hostname.includes('vercel.com') || hostname.includes('vercel.app')) {
    return { type: 'vercel', icon: 'cloud_upload' }
  }
  if (hostname.includes('netlify.com') || hostname.includes('netlify.app')) {
    return { type: 'netlify', icon: 'cloud_sync' }
  }
  if (hostname.includes('heroku.com') || hostname.includes('herokuapp.com')) {
    return { type: 'heroku', icon: 'cloud_queue' }
  }

  // Portfolio or personal website
  if (hostname.includes('portfolio') || hostname.includes('blog')) {
    return { type: 'portfolio', icon: 'work' }
  }

  // Default to website
  return { type: 'website', icon: 'language' }
}

export function isValidUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}
