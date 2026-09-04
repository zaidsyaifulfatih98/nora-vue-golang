export interface AnalyticsCountRow {
  label: string
  count: number
}

export interface AnalyticsSummary {
  totalPageViews: number
  totalCtaClicks: number
  todayPageViews: number
  todayCtaClicks: number
  ctaBreakdown: AnalyticsCountRow[]
  topPages: AnalyticsCountRow[]
}

export function useAnalyticsApi() {
  const axios = useAxios()
  const route = useRoute()

  // Fire-and-forget: tracking must never block or break the page a visitor
  // is actually using.
  const trackPageView = (path: string) =>
    axios.post('/analytics/track', { eventType: 'page_view', path }).catch(() => {})

  const trackCtaClick = (ctaLabel: string) =>
    axios.post('/analytics/track', { eventType: 'cta_click', path: route.path, ctaLabel }).catch(() => {})

  const getAnalyticsSummary = () =>
    axios.get('/analytics/summary').then((r) => r.data.data as AnalyticsSummary)

  return { trackPageView, trackCtaClick, getAnalyticsSummary }
}
