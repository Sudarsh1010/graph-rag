import { createFileRoute } from "@tanstack/react-router"
import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { MagnifyingGlassIcon } from "@phosphor-icons/react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "~/components/ui/card"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import { Label } from "~/components/ui/label"
import { apiFetch } from "~/lib/api-client"

function Search() {
  const [searchQuery, setSearchQuery] = useState("")
  const [submittedQuery, setSubmittedQuery] = useState("")

  const { data, isLoading, error } = useQuery({
    queryKey: ["search", submittedQuery],
    queryFn: () =>
      apiFetch<{ query: string; results: Record<string, Record<string, unknown>[]> }>(
        `/api/v1/search/?query=${encodeURIComponent(submittedQuery)}`,
      ),
    enabled: submittedQuery.length > 0,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setSubmittedQuery(searchQuery)
  }

  const resultGroups = data?.results ? Object.entries(data.results) : []

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="font-heading text-2xl font-bold tracking-tight">Search</h1>
        <p className="text-muted-foreground">
          Find products, orders, customers, and more across all entity types.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <MagnifyingGlassIcon className="size-5 text-primary" />
            Search
          </CardTitle>
          <CardDescription>Search across all SAP data entities by keyword.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="flex gap-2">
            <div className="flex-1">
              <Label htmlFor="search" className="sr-only">
                Search
              </Label>
              <Input
                id="search"
                type="search"
                placeholder="Type to search..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full"
              />
            </div>
            <Button type="submit" disabled={!searchQuery.trim()}>
              <MagnifyingGlassIcon className="size-4" />
              Search
            </Button>
          </form>
        </CardContent>
      </Card>

      {submittedQuery && (
        <div className="flex flex-col gap-4">
          <div className="text-sm text-muted-foreground">
            {isLoading
              ? "Searching..."
              : error
                ? "Error searching"
                : resultGroups.length > 0
                  ? `Found results across ${resultGroups.length} entity type(s)`
                  : `No results for "${submittedQuery}"`}
          </div>

          {isLoading && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              Searching...
            </div>
          )}

          {resultGroups.map(([entityType, results]) => (
            <Card key={entityType}>
              <CardHeader>
                <CardTitle className="capitalize text-base">
                  {entityType.replace(/-/g, " ")}
                </CardTitle>
                <CardDescription>{results.length} result(s)</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b">
                        {results.length > 0 &&
                          Object.keys(results[0])
                            .slice(0, 6)
                            .map((key) => (
                              <th key={key} className="px-3 py-2 text-left font-medium capitalize">
                                {key.replace(/_/g, " ")}
                              </th>
                            ))}
                      </tr>
                    </thead>
                    <tbody>
                      {results.slice(0, 20).map((row, i) => (
                        <tr key={i} className="border-b">
                          {Object.entries(row)
                            .slice(0, 6)
                            .map(([key, value]) => (
                              <td key={key} className="max-w-[200px] truncate px-3 py-2">
                                {value == null ? (
                                  <span className="text-muted-foreground">—</span>
                                ) : typeof value === "object" ? (
                                  <span className="text-xs text-muted-foreground">
                                    {JSON.stringify(value)}
                                  </span>
                                ) : (
                                  <span>{String(value)}</span>
                                )}
                              </td>
                            ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>
          ))}

          {!isLoading && resultGroups.length === 0 && !error && (
            <Card>
              <CardContent>
                <div className="flex flex-col items-center justify-center py-8">
                  <MagnifyingGlassIcon className="size-12 text-muted-foreground" />
                  <div className="mt-4 text-center">
                    <div className="font-medium">No results found</div>
                    <div className="text-sm text-muted-foreground">
                      Try adjusting your search terms
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      )}
    </div>
  )
}

export const Route = createFileRoute("/search/")({
  component: Search,
})
