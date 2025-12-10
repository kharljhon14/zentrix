import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';

export default function CompanyCard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Company Name</CardTitle>
      </CardHeader>
      <CardContent>
        <p>Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dicta, officia?</p>
      </CardContent>
      <CardFooter>
        <small>Insert related contacts</small>
      </CardFooter>
    </Card>
  );
}
