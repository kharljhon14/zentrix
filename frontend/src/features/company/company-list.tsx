import CompanyCard from './company-card';

export default function CompanyList() {
  return (
    <div className="grid grid-cols-4 gap-4 p-8">
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
      <CompanyCard />
    </div>
  );
}
