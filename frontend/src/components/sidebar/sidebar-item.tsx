import { SidebarMenuButton, SidebarMenuItem } from '../ui/sidebar';
import type { SidebarItem } from './sidebar';

interface Props {
  sidebarItem: SidebarItem;
}

export default function SidebarItem({ sidebarItem }: Props) {
  return (
    <SidebarMenuItem>
      <SidebarMenuButton asChild>
        <a
          href={sidebarItem.url}
          className="p-6"
        >
          {sidebarItem.icon && sidebarItem.icon}
          <span>{sidebarItem.name}</span>
        </a>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
}
