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
import {
  Zap,
  Shield,
  Globe,
  Code,
  DollarSign,
  Cpu,
  SquareStack,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { AnimateInView } from '@/components/animate-in-view'

interface FeaturesProps {
  className?: string
}

export function Features(_props: FeaturesProps) {
  const { t } = useTranslation()

  const features = [
    {
      id: 'models',
      num: '01',
      title: t('DeepSeek · Claude · Qwen'),
      desc: t('支持 DeepSeek V4 Flash/Pro、Claude Sonnet/Haiku、通义千问 Max/Plus 等中国主流模型，一个 API 统一调用'),
      span: 'md:col-span-2',
      icon: <Cpu className='size-4 text-violet-400' />,
      visual: (
        <div className='mt-4 grid grid-cols-3 gap-2'>
          {['DeepSeek', 'Claude', 'Qwen', 'GLM', 'Doubao', 'Spark'].map(
            (name) => (
              <div
                key={name}
                className='border-border/30 bg-muted/20 text-muted-foreground flex items-center justify-center rounded-lg border px-3 py-2 text-xs transition-colors duration-300 hover:border-violet-500/30 hover:bg-violet-500/5'
              >
                {name}
              </div>
            )
          )}
        </div>
      ),
    },
    {
      id: 'price',
      num: '02',
      title: t('比官方更低价格'),
      desc: t('通过集中采购和线路优化，价格比官方渠道便宜 10-20%，按量计费无最低消费'),
      span: 'md:col-span-1',
      icon: <DollarSign className='size-4 text-emerald-400' />,
      visual: (
        <div className='mt-4 flex items-center justify-center'>
          <div className='relative'>
            <div className='flex size-16 items-center justify-center rounded-2xl border border-emerald-500/20 bg-emerald-500/5'>
              <DollarSign
                className='size-7 text-emerald-500/70'
                strokeWidth={1.5}
              />
            </div>
          </div>
        </div>
      ),
    },
    {
      id: 'global',
      num: '03',
      title: t('全球加速'),
      desc: t('美西服务器，优化的跨境线路。DeepSeek 仅 300ms 响应'),
      span: 'md:col-span-1',
      icon: <Globe className='size-4 text-blue-400' />,
      visual: (
        <div className='mt-4 space-y-2'>
          {[t('全球节点'), t('智能路由'), t('自动切换')].map(
            (step, i) => (
              <div key={step} className='flex items-center gap-2'>
                <div
                  className={`flex size-6 items-center justify-center rounded-full text-[10px] font-bold ${
                    i === 1
                      ? 'border border-blue-500/30 bg-blue-500/20 text-blue-500'
                      : 'border-border/40 bg-muted text-muted-foreground border'
                  }`}
                >
                  {i + 1}
                </div>
                <div className='bg-border/40 h-px flex-1' />
                <span className='text-muted-foreground text-xs'>{step}</span>
              </div>
            )
          )}
        </div>
      ),
    },
    {
      id: 'developer',
      num: '04',
      title: t('双支付系统'),
      desc: t('支持支付宝和 Stripe 支付，美元结算，全球开发者无门槛接入'),
      span: 'md:col-span-2',
      icon: <SquareStack className='size-4 text-amber-400' />,
      visual: (
        <div className='mt-4 flex items-center gap-3'>
          <div className='flex -space-x-2'>
            {['Alipay', 'Stripe', 'USD', 'API'].map((n) => (
              <div
                key={n}
                className='border-background from-muted to-muted/60 text-muted-foreground flex size-8 items-center justify-center rounded-full border-2 bg-gradient-to-br text-[9px] font-bold'
              >
                {n}
              </div>
            ))}
          </div>
          <div className='text-muted-foreground flex items-center gap-1.5 text-xs'>
            <Code className='size-3.5 text-blue-500' />
            {t('OpenAI 兼容')}
          </div>
        </div>
      ),
    },
  ]

  const additionalFeatures = [
    {
      icon: <Zap className='size-5' strokeWidth={1.5} />,
      title: t('极速响应'),
      desc: t('中美优化线路，Token 生成速度极快'),
    },
    {
      icon: <Shield className='size-5' strokeWidth={1.5} />,
      title: t('安全可靠'),
      desc: t('企业级安全，完善的权限管理和用量控制'),
    },
    {
      icon: <Globe className='size-5' strokeWidth={1.5} />,
      title: t('无需实名'),
      desc: t('不需要中国手机号验证，注册即用'),
    },
    {
      icon: <Code className='size-5' strokeWidth={1.5} />,
      title: t('零迁移成本'),
      desc: t('兼容 OpenAI SDK，一行代码切换 base_url'),
    },
  ]

  return (
    <section className='relative z-10 px-6 py-24 md:py-32'>
      <div className='mx-auto max-w-6xl'>
        <AnimateInView className='mb-16 max-w-lg'>
          <p className='text-muted-foreground mb-3 text-xs font-medium tracking-widest uppercase'>
            {t('核心优势')}
          </p>
          <h2 className='text-2xl leading-tight font-bold tracking-tight md:text-3xl'>
            {t('为什么选择 ChinaiAPI？')}
          </h2>
        </AnimateInView>

        {/* Bento grid */}
        <div className='border-border/40 bg-border/40 grid gap-px overflow-hidden rounded-xl border md:grid-cols-3'>
          {features.map((f, i) => (
            <AnimateInView
              key={f.id}
              delay={i * 100}
              animation='scale-in'
              className={`bg-background group hover:bg-muted/20 p-7 transition-colors duration-300 md:p-8 ${f.span}`}
            >
              <div className='mb-3 flex items-center gap-3'>
                <span className='border-border/40 bg-muted text-muted-foreground flex size-7 items-center justify-center rounded-md border text-[10px] font-semibold tabular-nums'>
                  {f.num}
                </span>
                <h3 className='text-sm font-semibold'>{f.title}</h3>
              </div>
              <p className='text-muted-foreground text-sm leading-relaxed'>
                {f.desc}
              </p>
              {f.visual}
            </AnimateInView>
          ))}
        </div>

        {/* Additional features row */}
        <div className='mt-12 grid grid-cols-2 gap-8 md:grid-cols-4 md:gap-12'>
          {additionalFeatures.map((f, i) => (
            <AnimateInView
              key={f.title}
              delay={i * 100}
              animation='fade-up'
              className='flex flex-col items-center text-center'
            >
              <div className='text-muted-foreground border-border/50 bg-muted/30 group-hover:text-foreground mb-3 flex size-12 items-center justify-center rounded-xl border transition-colors'>
                {f.icon}
              </div>
              <h3 className='mb-1.5 text-sm font-semibold'>{f.title}</h3>
              <p className='text-muted-foreground max-w-[200px] text-xs leading-relaxed'>
                {f.desc}
              </p>
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}
