import { createFileRoute, Outlet } from '@tanstack/react-router';
import React from 'react';

import Header from '@/components/header/header';
import AppSidebar from '@/components/sidebar/sidebar';
import AppSidebarProvider from '@/providers/sidebar-provider';

export const Route = createFileRoute('/_dashboard')({
  component: RouteComponent
});

function RouteComponent() {
  return (
    <React.Fragment>
      <AppSidebarProvider>
        <AppSidebar />
        <div className="w-full">
          <Header />

          <main className="">
            <Outlet />
          </main>
        </div>
      </AppSidebarProvider>
    </React.Fragment>
  );
}
