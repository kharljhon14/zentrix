import { SidebarProvider } from '@/components/ui/sidebar';
import type { PropsWithChildren } from 'react';

export default function AppSidebarProvider({ children }: PropsWithChildren) {
  return <SidebarProvider>{children}</SidebarProvider>;
}
