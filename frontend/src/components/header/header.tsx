import { PanelLeft } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Button } from '../ui/button';
import { useSidebar } from '../ui/sidebar';
import { Separator } from '../ui/separator';

export default function Header() {
  const { toggleSidebar } = useSidebar();
  return (
    <header className="shadow h-16 p-4">
      <div className="w-full h-full flex items-center justify-between">
        <Button
          variant="ghost"
          onClick={toggleSidebar}
        >
          <PanelLeft />
        </Button>
        <div className="ml-2 h-full mr-auto">
          <Separator orientation="vertical" />
        </div>
        <Avatar>
          <AvatarImage src="https://github.com/shadcn.png" />
          <AvatarFallback>CN</AvatarFallback>
        </Avatar>
      </div>
    </header>
  );
}
