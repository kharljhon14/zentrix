import { CirclePlus } from 'lucide-react';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog';

import CompanyForm from './company-form';

export default function NewCompanyModal() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>
          <CirclePlus />
          Add New Company
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Company</DialogTitle>
        </DialogHeader>
        <CompanyForm />
        <DialogFooter>
          <div>
            <Button
              variant="outline"
              className="mr-2"
            >
              Cancel
            </Button>
            <Button>Save</Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
