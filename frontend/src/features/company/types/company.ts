export interface Company {
  id: string;
  name: string;
  address: string;
  sales_owner?: string;
  sales_owner_name?: string;
  email: string;
  company_size: string;
  industry: string;
  business_type: string;
  country: string;
  image?: string;
  website: string;
  created_at: string;
  updated_at: string;
}
