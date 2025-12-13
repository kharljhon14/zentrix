import { createFileRoute } from '@tanstack/react-router';

import CompanyList from '@/features/company/components/company-list';
import NewCompanyModal from '@/features/company/components/new-company-modal';

export const Route = createFileRoute('/_dashboard/companies')({
  component: RouteComponent
});

function RouteComponent() {
  return (
    <div className="p-8">
      <div className="ml-auto w-fit mb-6">
        <NewCompanyModal />
      </div>
      <CompanyList />
    </div>
  );
}
