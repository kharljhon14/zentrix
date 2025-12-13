import { useQuery } from '@tanstack/react-query';

import agent from '@/api/agent';

export default function useGetCompany(id: string) {
  const query = useQuery({
    queryKey: [id],
    queryFn: () => agent.companies.getById(id)
  });

  return query;
}
