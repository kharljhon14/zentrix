import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { Link } from '@tanstack/react-router';
import { Ellipsis } from 'lucide-react';

export default function CompanyCard() {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Company Name</CardTitle>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                size="icon-sm"
                variant="ghost"
              >
                <Ellipsis />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <div className="flex flex-col gap-2">
                <Link to="/">View Company</Link>
                <button className=" cursor-pointer">Delete Company</button>
              </div>
            </TooltipContent>
          </Tooltip>
        </div>
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
