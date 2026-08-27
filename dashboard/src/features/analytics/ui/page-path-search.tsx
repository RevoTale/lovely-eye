import { Search } from 'lucide-react';
import { type FormEvent, type FunctionComponent, useEffect, useState } from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';

interface PagePathSearchProps {
  value: string;
  onSearch: (value: string) => void;
}

const MAX_PATH_LENGTH = 2048;

const PagePathSearch: FunctionComponent<PagePathSearchProps> = ({ value, onSearch }) => {
  const [draft, setDraft] = useState(value);

  useEffect(() => setDraft(value), [value]);

  const submit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    onSearch(draft);
  };

  return (
    <form
      className='mb-4 flex flex-col gap-2 xs:flex-row'
      aria-label='Search top pages by path'
      onSubmit={submit}
    >
      <Input
        type='search'
        value={draft}
        maxLength={MAX_PATH_LENGTH}
        placeholder='/blog'
        aria-label='Page path contains'
        onChange={(event) => setDraft(event.target.value)}
      />
      <Button type='submit' variant='secondary' className='w-full xs:w-auto'>
        <Search data-icon='inline-start' />
        Search
      </Button>
    </form>
  );
};

export default PagePathSearch;
