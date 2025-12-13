import { useQuery } from '@tanstack/react-query';

import agent from '@/api/agent';

export default function useGetCompanies() {
  const query = useQuery({
    queryKey: ['companies'],
    queryFn: () => agent.companies.getCompanies()
  });

  return query;
}
