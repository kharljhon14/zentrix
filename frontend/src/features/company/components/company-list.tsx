import { Card } from '@/components/ui/card';
import useGetCompanies from '../hooks/api/use-get-companies';
import CompanyCard from './company-card';
import { Skeleton } from '@/components/ui/skeleton';

export default function CompanyList() {
  const companiesQuery = useGetCompanies();

  if (companiesQuery.isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <Skeleton className="h-44 w-64 rounded-xl bg-neutral-400" />
        <Skeleton className="h-44 w-64 rounded-xl bg-neutral-400" />
        <Skeleton className="h-44 w-64 rounded-xl bg-neutral-400" />
        <Skeleton className="h-44 w-64 rounded-xl bg-neutral-400" />
        <Skeleton className="h-44 w-64 rounded-xl bg-neutral-400" />
      </div>
    );
  }

  if (companiesQuery.isSuccess)
    return <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6"></div>;
}
