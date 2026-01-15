import React from 'react';
import { Globe } from 'lucide-react';
import { Badge, Card, CardContent, CardHeader, CardTitle, Progress } from '@/components/ui';
import { Link } from '@/router';
import type { CountryStats } from '@/generated/graphql';
import { addFilterValue } from '@/lib/filter-utils';

interface CountryCardProps {
  countries: CountryStats[];
  siteId: string;
}

export function CountryCard({ countries, siteId }: CountryCardProps): React.JSX.Element {
  const totalVisitors = countries.reduce((sum, item) => sum + item.visitors, 0);

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <Globe className="h-4 w-4 text-primary" />
          </div>
          Countries
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {countries.length > 0 ? (
            countries.map((countryStat, index) => {
              const percentage = totalVisitors > 0 ? (countryStat.visitors / totalVisitors) * 100 : 0;

              return (
                <div key={index}>
                  <div className="flex items-center justify-between mb-1">
                    <Link
                      to="/sites/$siteId"
                      params={{ siteId }}
                      search={(prev) => ({
                        ...prev,
                        country: addFilterValue(prev.country, countryStat.country),
                      })}
                      className="text-sm font-medium truncate max-w-[200px] hover:text-primary hover:underline cursor-pointer"
                    >
                      {countryStat.country}
                    </Link>
                    <div className="flex items-center gap-2">
                      <Badge variant="secondary">{countryStat.visitors.toLocaleString()}</Badge>
                      <span className="text-xs text-muted-foreground">{percentage.toFixed(1)}%</span>
                    </div>
                  </div>
                  <Progress value={percentage} className="h-2" />
                </div>
              );
            })
          ) : (
            <p className="text-sm text-muted-foreground text-center py-4">No country data yet</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
