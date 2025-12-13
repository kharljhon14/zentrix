import type { PropsWithChildren } from 'react';

import { SidebarProvider } from '@/components/ui/sidebar';

export default function AppSidebarProvider({ children }: PropsWithChildren) {
  return <SidebarProvider>{children}</SidebarProvider>;
}
