import { ChevronsUpDown } from 'lucide-react';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';
import { SidebarMenuButton, SidebarMenuItem } from '../ui/sidebar';
import type { SidebarItem } from './sidebar';

interface Props {
  sidebarItem: SidebarItem;
}

export default function SidebarItem({ sidebarItem }: Props) {
  return (
    <SidebarMenuItem>
      {sidebarItem.subItems ? (
        <Collapsible>
          <CollapsibleTrigger className="cursor-pointer flex gap-2 px-6 py-4  w-full hover:bg-gray-100 rounded-md">
            {sidebarItem.icon && sidebarItem.icon}
            <span className="flex justify-between items-center w-full">
              {sidebarItem.name} <ChevronsUpDown size={16} />
            </span>
          </CollapsibleTrigger>
          <CollapsibleContent className="pl-6">
            {sidebarItem.subItems.map((subitem) => (
              <SidebarItem
                key={subitem.name}
                sidebarItem={subitem}
              />
            ))}
          </CollapsibleContent>
        </Collapsible>
      ) : (
        <SidebarMenuButton asChild>
          <a
            href={sidebarItem.url}
            className="p-6"
          >
            {sidebarItem.icon && sidebarItem.icon}
            <span>{sidebarItem.name}</span>
          </a>
        </SidebarMenuButton>
      )}
    </SidebarMenuItem>
  );
}
