import { createFileRoute } from "@tanstack/react-router"
import { useState } from "react"
import { useMutation } from "@tanstack/react-query"
import { GraphIcon } from "@phosphor-icons/react"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card"
import { Button } from "~/components/ui/button"
import { Textarea } from "~/components/ui/textarea"
import { Label } from "~/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table"
import { apiFetch } from "~/lib/api-client"
import { Skeleton } from "~/components/ui/skeleton"

const EXAMPLE_QUERIES = [
  {
    label: "All customers",
    cypher: `MATCH (c:Customer) RETURN c LIMIT 20`,
  },
  {
    label: "Orders and their items",
    cypher: `MATCH (so:SalesOrder)-[:CONTAINS_ITEM]->(item:SalesOrderItem) RETURN so, item LIMIT 10`,
  },
  {
    label: "Order-to-delivery chain",
    cypher: `MATCH (so:SalesOrder)-[:CONTAINS_ITEM]->(item:SalesOrderItem)-[:REFERENCES_PRODUCT]->(p:Product) RETURN so.salesOrder, item.material, p.product LIMIT 10`,
  },
  {
    label: "Delivery to billing",
    cypher: `MATCH (d:Delivery) RETURN d LIMIT 10`,
  },
]

function GraphView() {
  const [cypher, setCypher] = useState("")
  const [result, setResult] = useState<{ data: Array<unknown> } | null>(null)
  const [error, setError] = useState<string | null>(null)

  const mutation = useMutation({
    mutationFn: (query: string) =>
      apiFetch<{ data: Array<unknown> }>("/api/v1/graph/query", {
        method: "POST",
        body: JSON.stringify({ cypher: query }),
      }),
    onSuccess: (data) => {
      setResult(data)
      setError(null)
    },
    onError: (err) => {
      setError(err instanceof Error ? err.message : "Query failed")
      setResult(null)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!cypher.trim()) return
    mutation.mutate(cypher)
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="font-heading text-2xl font-bold tracking-tight">
          Graph Explorer
        </h1>
        <p className="text-muted-foreground">
          Query the SAP data graph using Apache AGE Cypher queries.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <GraphIcon className="size-5 text-primary" />
            Cypher Query
          </CardTitle>
          <CardDescription>
            Execute read-only Cypher queries against the SAP order-to-cash
            graph.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="flex flex-col gap-3">
            <Label htmlFor="cypher">Cypher Query</Label>
            <Textarea
              id="cypher"
              placeholder="MATCH (n) RETURN n LIMIT 10"
              value={cypher}
              onChange={(e) => setCypher(e.target.value)}
              className="min-h-[100px] font-mono text-sm"
            />
            <div className="flex items-center justify-between">
              <div className="text-xs text-muted-foreground">
                {cypher.length} characters
              </div>
              <Button
                type="submit"
                disabled={!cypher.trim() || mutation.isPending}
              >
                <GraphIcon className="size-4" />
                {mutation.isPending ? "Executing..." : "Run Query"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Example Queries</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {EXAMPLE_QUERIES.map((eq) => (
              <Button
                key={eq.label}
                variant="outline"
                size="sm"
                onClick={() => setCypher(eq.cypher)}
              >
                {eq.label}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {error && (
        <Card className="border-destructive">
          <CardContent className="pt-6">
            <div className="text-sm text-destructive">{error}</div>
          </CardContent>
        </Card>
      )}

      {mutation.isPending && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Results</CardTitle>
            <CardDescription>
              <Skeleton className="h-4 w-24" />
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    {[1, 2, 3, 4].map((i) => (
                      <TableHead key={i}>
                        <Skeleton className="h-4 w-20" />
                      </TableHead>
                    ))}
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      {[1, 2, 3, 4].map((j) => (
                        <TableCell key={j}>
                          <Skeleton className="h-4 w-full" />
                        </TableCell>
                      ))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      )}

      {result && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Results</CardTitle>
            <CardDescription>
              {Array.isArray(result.data)
                ? `${result.data.length} row(s) returned`
                : "No data"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {Array.isArray(result.data) && result.data.length > 0 ? (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      {Object.keys(
                        result.data[0] as Record<string, unknown>
                      ).map((key) => (
                        <TableHead key={key}>
                          {key}
                        </TableHead>
                      ))}
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {result.data.slice(0, 50).map((row, i) => (
                      <TableRow key={i}>
                        {Object.entries(row as Record<string, unknown>).map(
                          ([key, value]) => (
                            <TableCell
                              key={key}
                              className="max-w-[300px] truncate"
                            >
                              {typeof value === "object" ? (
                                <span className="font-mono text-xs">
                                  {JSON.stringify(value)}
                                </span>
                              ) : (
                                <span>{String(value ?? "—")}</span>
                              )}
                            </TableCell>
                          )
                        )}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            ) : (
              <div className="py-8 text-center text-muted-foreground">
                No results returned
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export const Route = createFileRoute("/graph/")({
  component: GraphView,
})
