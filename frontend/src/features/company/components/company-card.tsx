import { Link } from '@tanstack/react-router';
import { Ellipsis } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';

import type { Company } from '../types/company';

interface Props {
  company: Company;
}

export default function CompanyCard({ company }: Props) {
  return (
    <Card>
      <CardHeader className="relative">
        <div className="flex items-center justify-center">
          <CardTitle>
            <div className="flex flex-col items-center gap-2 justify-center">
              <div className="rounded-2xl w-16 h-16 overflow-hidden block">
                <img
                  className="w-full h-full object-cover"
                  src={company.image ? company.image : 'https://placehold.net/building-400x400.png'}
                  alt={company.name}
                />
              </div>
              <p>{company.name}</p>
            </div>
          </CardTitle>
          <Popover>
            <PopoverTrigger
              asChild
              className="absolute -top-2 right-4"
            >
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
        <div className="flex justify-between w-full">
          <div>
            <small>Related contacts</small>
          </div>
          <div className="text-xs">
            <p>Sales owner</p>
            <p>{company.sales_owner_name}</p>
          </div>
        </div>
      </CardFooter>
    </Card>
  );
}
