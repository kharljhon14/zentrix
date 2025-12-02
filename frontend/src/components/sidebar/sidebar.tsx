import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu
} from '@/components/ui/sidebar';
import {
  BookUser,
  Building2,
  Calendar,
  Gauge,
  NotepadText,
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
    icon: <SquareKanban size={18} />,
    name: 'Scrumboard',
    subItems: [
      {
        name: 'Project Kanban',
        url: '/kaban'
      },
      {
        name: 'Sales Pipeline',
        url: '/sales'
      }
    ]
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
    icon: <ShieldUser size={18} />,
    name: 'Administration',
    subItems: [
      {
        name: 'Settings',
        url: '/settings'
      },
      {
        name: 'Audit Log',
        url: '/logs'
      }
    ]
  }
];

export default function AppSidebar() {
  return (
    <Sidebar>
      <SidebarContent className="bg-white">
        <SidebarGroup>
          <SidebarGroupLabel>Tenant Name</SidebarGroupLabel>
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
