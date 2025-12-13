import { useForm } from 'react-hook-form';

import { Form, FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form';
import { Input } from '@/components/ui/input';

export default function CompanyForm() {
  const form = useForm();

  return (
    <Form {...form}>
      <form
        autoComplete="off"
        className="space-y-4 pt-4"
      >
        <FormField
          name="name"
          control={form.control}
          render={({ field }) => (
            <FormItem>
              <FormLabel>Company Name</FormLabel>
              <FormControl>
                <Input
                  placeholder="Please enter company name"
                  value={field.value}
                  onChange={field.onChange}
                />
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          name="sales_owner"
          control={form.control}
          render={({ field }) => (
            <FormItem>
              <FormLabel>Sales Owner</FormLabel>
              <FormControl>
                <Input
                  placeholder="Please enter sales owner"
                  value={field.value}
                  onChange={field.onChange}
                />
              </FormControl>
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
