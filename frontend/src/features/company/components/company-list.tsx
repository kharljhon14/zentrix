import useGetCompany from '../hooks/api/use-get-company';
import CompanyCard from './company-card';

export default function CompanyList() {
  const companyQuery = useGetCompany('7bdcc34e-4236-4ef0-8efb-eda5b915e6a9');

  console.table(companyQuery.data);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
    </div>
  );
}
