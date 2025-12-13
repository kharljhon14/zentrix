import { Link } from '@tanstack/react-router';
import { Ellipsis } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';

export default function CompanyCard() {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Company Name</CardTitle>
          <Popover>
            <PopoverTrigger asChild>
              <Button
                size="icon-sm"
                variant="ghost"
              >
                <Ellipsis />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-40">
              <div className="flex flex-col gap-2">
                <Button
                  size="sm"
                  asChild
                  variant="outline"
                >
                  <Link to="/">View Company</Link>
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                >
                  Delete Company
                </Button>
              </div>
            </PopoverContent>
          </Popover>
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
