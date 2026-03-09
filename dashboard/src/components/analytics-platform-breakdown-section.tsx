import type { FunctionComponent } from 'react';
import BrowserCard from '@/components/browser-card';
import DevicesCard from '@/components/devices-card';
import OSCard from '@/components/os-card';
import type { BrowserStatsFieldsFragment, DeviceStatsFieldsFragment, OperatingSystemStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface AnalyticsPlatformBreakdownSectionProps {
  siteId: string;
  dashboardState: DashboardLoadState;
  totalVisitors: number;
  browsers: BrowserStatsFieldsFragment[];
  devices: DeviceStatsFieldsFragment[];
  devicesTotal: number;
  devicesTotalVisitors: number;
  devicesPage: number;
  devicesPageSize: number;
  onDevicesPageChange: (page: number) => void;
  operatingSystems: OperatingSystemStatsFieldsFragment[];
  operatingSystemsTotal: number;
  operatingSystemsTotalVisitors: number;
  osPage: number;
  osPageSize: number;
  onOSPageChange: (page: number) => void;
}

const AnalyticsPlatformBreakdownSection: FunctionComponent<AnalyticsPlatformBreakdownSectionProps> = (props) => (
  <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
    <BrowserCard browsers={props.browsers} totalVisitors={props.totalVisitors} siteId={props.siteId} state={props.dashboardState} />
    <DevicesCard devices={props.devices} total={props.devicesTotal} totalVisitors={props.devicesTotalVisitors} page={props.devicesPage} pageSize={props.devicesPageSize} siteId={props.siteId} state={props.dashboardState} onPageChange={props.onDevicesPageChange} />
    <OSCard operatingSystems={props.operatingSystems} total={props.operatingSystemsTotal} totalVisitors={props.operatingSystemsTotalVisitors} page={props.osPage} pageSize={props.osPageSize} siteId={props.siteId} state={props.dashboardState} onPageChange={props.onOSPageChange} />
  </div>
);

export default AnalyticsPlatformBreakdownSection;
