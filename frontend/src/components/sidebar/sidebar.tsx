import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from '@/components/ui/sidebar';
import {
  BookUser,
  Building2,
  Calendar,
  Gauge,
  LayoutList,
  Logs,
  NotepadText,
  Settings,
  ShieldUser,
  SquareKanban
} from 'lucide-react';
import type { JSX } from 'react';
import SidebarItem from './sidebar-item';

export type SidebarItem = {
  icon?: JSX.Element;
  name: string;
  url?: string;
  subItems?: SidebarItem[];
};

const sidebarItems: SidebarItem[] = [
  {
    icon: <Gauge />,
    name: 'Dashboard',
    url: '/dashboard'
  },
  {
    icon: <Calendar />,
    name: 'Calendar',
    url: '/calendar'
  },

  {
    icon: <LayoutList />,
    name: 'Projects',
    url: '/kaban'
  },
  {
    icon: <SquareKanban />,
    name: 'Sales Pipeline',
    url: '/sales'
  },

  {
    icon: <Building2 />,
    name: 'Companies',
    url: '/companies'
  },
  {
    icon: <BookUser />,
    name: 'Contacts',
    url: '/contacts'
  },
  {
    icon: <NotepadText />,
    name: 'Quotes',
    url: '/quotes'
  },
  {
    icon: <Logs />,
    name: 'Audit Log',
    url: '/logs'
  },
  {
    icon: <Settings />,
    name: 'Settings',
    url: '/settings'
  }
];

export default function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild>
              <a
                href="#"
                className="p-6"
              >
                <span className="text-lg font-semibold">Acme Inc.</span>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent className="bg-white">
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {sidebarItems.map((item) => (
                <SidebarItem
                  key={item.name}
                  sidebarItem={item}
                />
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
