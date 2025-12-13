import { Skeleton } from '@/components/ui/skeleton';

import useGetCompanies from '../hooks/api/use-get-companies';
import CompanyCard from './company-card';

export default function CompanyList() {
  const companiesQuery = useGetCompanies();

  if (companiesQuery.isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
        <Skeleton className="h-64 w-80 rounded-xl bg-neutral-400" />
      </div>
    );
  }

  if (companiesQuery.isSuccess)
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {companiesQuery.data.data.map((company) => (
          <CompanyCard
            key={company.id}
            company={company}
          />
        ))}
      </div>
    );
}
