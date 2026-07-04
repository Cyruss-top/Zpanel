import type { GlobalThemeOverrides } from 'naive-ui'

export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#2563EB',
    primaryColorHover: '#1D4ED8',
    primaryColorPressed: '#1E40AF',
    primaryColorSuppl: '#3B82F6',
    borderRadius: '4px',
    borderRadiusSmall: '4px',
    fontFamily: '"Source Sans 3", "PingFang SC", "Microsoft YaHei", system-ui, sans-serif',
    fontFamilyMono: '"IBM Plex Mono", "SF Mono", Consolas, monospace',
    textColorBase: '#18181B',
    textColor1: '#18181B',
    textColor2: '#71717A',
    textColor3: '#A1A1AA',
    bodyColor: '#F4F4F5',
    dividerColor: '#E4E4E7',
    borderColor: '#E4E4E7',
  },
  Layout: {
    color: '#F4F4F5',
    siderColor: '#18181B',
    headerColor: '#FFFFFF',
    footerColor: '#F4F4F5',
  },
  Menu: {
    itemTextColor: '#A1A1AA',
    itemTextColorHover: '#FAFAFA',
    itemTextColorActive: '#FAFAFA',
    itemColorActive: '#27272A',
    itemColorHover: '#27272A',
    borderRadius: '4px',
  },
  Card: {
    color: '#FFFFFF',
    borderColor: '#E4E4E7',
    borderRadius: '4px',
  },
  Button: {
    borderRadiusSmall: '4px',
    borderRadiusMedium: '4px',
    borderRadiusLarge: '4px',
  },
  DataTable: {
    borderRadius: '4px',
    thColor: '#FAFAFA',
    tdColor: '#FFFFFF',
  },
  Tag: {
    borderRadius: '4px',
  },
}
