/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
/**
 * Home page constants
 * All hardcoded data for home page sections
 */
import { type TFunction } from 'i18next'

// Layout - Main base classes
export const MAIN_BASE_CLASSES = 'bg-background text-foreground w-full'

// Hero section - AI Applications (Left side)
export const AI_APPLICATIONS = [
  'LobeHub.Color',
  'Dify.Color',
  'OpenWebUI',
  'Cline',
] as const

// Hero section - AI Models (Right side)
export const AI_MODELS = [
  'Qwen.Color',
  'DeepSeek.Color',
  'Doubao.Color',
  'OpenAI',
  'Claude.Color',
  'Gemini.Color',
] as const

// Hero section - Gateway Features
export const GATEWAY_FEATURES = [
  'Cost Tracking',
  'Model Access',
  'Guardrails',
  'Observability',
  'Budgets',
  'Load Balancing',
  'Rate Limiting',
  'Token Mgmt',
  'Prompt Caching',
  'Pass-Through',
] as const

// Stats section - Default statistics
export const DEFAULT_STATS = [
  {
    value: '40',
    suffix: '+',
    description: 'AI Providers',
  },
  {
    value: '200',
    suffix: '+',
    description: 'AI Models',
  },
  {
    value: 'OpenAI',
    suffix: '',
    description: 'OpenAI Compatible',
  },
  {
    value: '99.9',
    suffix: '%',
    description: 'Uptime SLA',
  },
] as const

// Features section - Default features
export const DEFAULT_FEATURES = [
  {
    title: 'Universal API Access',
    description: 'One standard endpoint for all major AI models',
    iconName: 'Zap',
  },
  {
    title: 'Enterprise Security',
    description: 'Role-based access control, SSO, and key management',
    iconName: 'Shield',
  },
  {
    title: 'Intelligent Routing',
    description: 'Auto load balancing, failover, and cost optimization',
    iconName: 'Globe',
  },
  {
    title: 'Developer First',
    description: 'Drop-in replacement for OpenAI SDK with comprehensive docs',
    iconName: 'Code',
  },
  {
    title: 'Usage Analytics',
    description: 'Real-time dashboards, cost tracking, and token monitoring',
    iconName: 'Gauge',
  },
  {
    title: 'Flexible Billing',
    description: 'Pay-as-you-go, prepaid credits, and subscription plans',
    iconName: 'DollarSign',
  },
  {
    title: 'Multi-Tenant',
    description: 'Isolated workspaces with per-team quotas and access policies',
    iconName: 'Users',
  },
  {
    title: 'Easy Integration',
    description: 'Docker deployment, environment-driven, minimal configuration',
    iconName: 'Download',
  },
] as const

export function getGatewayFeatures(t: TFunction) {
  return GATEWAY_FEATURES.map((feature) => t(feature))
}

export function getDefaultStats(t: TFunction) {
  return DEFAULT_STATS.map((stat) => ({
    ...stat,
    description: stat.description ? t(stat.description) : undefined,
  }))
}

export function getDefaultFeatures(t: TFunction) {
  return DEFAULT_FEATURES.map((feature) => ({
    ...feature,
    title: t(feature.title),
    description: t(feature.description),
  }))
}
