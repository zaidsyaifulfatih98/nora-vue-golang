// Shared open/closed state for the mobile drawer sidebar, so the hamburger
// button in DashboardHeader.vue and the drawer in SideBar.vue can toggle
// the same state without prop-drilling through the dashboard layout.
export function useDashboardSidebar() {
  return useState('dashboard-sidebar-open', () => false)
}
