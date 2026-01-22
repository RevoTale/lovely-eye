
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import { FileText, Users } from 'lucide-react';

interface ActivePage {
  path: string;
  visitors: number;
}

interface ActivePagesCardProps {
  activePages: ActivePage[];
}

const EMPTY_COUNT = 0;

export function ActivePagesCard({ activePages }: ActivePagesCardProps): React.JSX.Element {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
              <FileText className="h-4 w-4 text-primary" />
            </div>
            Active Pages
          </div>
          <Badge variant="outline" className="flex items-center gap-1">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            Live
          </Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {activePages.length > EMPTY_COUNT ? (
            activePages.map((page, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 rounded-lg bg-muted/50 hover:bg-muted transition-colors"
              >
                <div className="flex items-center gap-3 flex-1 min-w-0">
                  <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                    <FileText className="h-5 w-5 text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{page.path}</p>
                    <p className="text-xs text-muted-foreground">
                      Currently viewing
                    </p>
                  </div>
                </div>
                <Badge variant="secondary" className="flex items-center gap-1 ml-2">
                  <Users className="h-3 w-3" />
                  {page.visitors}
                </Badge>
              </div>
            ))
          ) : (
            <div className="text-center py-8">
              <FileText className="h-12 w-12 mx-auto mb-3 text-muted-foreground opacity-50" />
              <p className="text-sm text-muted-foreground">No active visitors</p>
              <p className="text-xs text-muted-foreground mt-1">
                Pages will appear as visitors browse your site
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
