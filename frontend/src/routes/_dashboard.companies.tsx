import CompanyList from '@/features/company/company-list';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/companies')({
  component: RouteComponent
});

function RouteComponent() {
  return (
    <div>
      <CompanyList />
    </div>
  );
}
