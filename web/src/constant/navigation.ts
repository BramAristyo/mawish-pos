import { LayoutDashboard, Package, Receipt, Users, Wallet } from 'lucide-vue-next'

export const NAVIGATION_CONFIG = [
  {
    title: 'Overview',
    url: '#',
    icon: LayoutDashboard,
    items: [
      {
        title: 'Dashboard',
        url: '/dashboard',
      },
    ],
  },
  {
    title: 'Menu',
    url: '#',
    icon: Package,
    items: [
      {
        title: 'Products',
        url: '/catalog/products',
      },
      {
        title: 'Categories',
        url: '/categories',
      },
      {
        title: 'Modifiers',
        url: '/modifiers',
      },
    ],
  },
  {
    title: 'Transactions',
    url: '#',
    icon: Receipt,
    items: [
      {
        title: 'Taxes',
        url: '/taxes',
      },
      {
        title: 'Sales Types',
        url: '/sales-types',
      },
      {
        title: 'Discounts',
        url: '/discounts',
      },
    ],
  },
  {
    title: 'Accounting',
    url: '#',
    icon: Wallet,
    items: [
      {
        title: 'Chart of Accounts',
        url: '/coa',
      },
      {
        title: 'Input Ledger',
        url: '/accounting/cash-transactions',
      },
    ],
  },
  {
    title: 'People',
    url: '#',
    icon: Users,
    items: [
      {
        title: 'Employee',
        url: '/employees',
      },
      {
        title: 'Shift Schedule',
        url: '/shift-schedules',
      },
      {
        title: 'Attendance',
        url: '/attendances',
      },
      {
        title: 'Payroll',
        url: '/payroll',
      },
    ],
  },
]
