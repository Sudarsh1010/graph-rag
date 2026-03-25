import { createFileRoute } from "@tanstack/react-router"
import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import {
  CreditCardIcon,
  FileTextIcon,
  PackageIcon,
  ShoppingCartIcon,
  TruckIcon,
  UsersIcon,
} from "@phosphor-icons/react"
import type { PaginatedResponse } from "~/lib/api-types"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card"
import { Button } from "~/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table"
import { apiFetch } from "~/lib/api-client"

type EntityType =
  | "sales-orders"
  | "products"
  | "customers"
  | "deliveries"
  | "billing-documents"
  | "payments"

const entityTypes = [
  {
    id: "sales-orders" as EntityType,
    name: "Sales Orders",
    icon: ShoppingCartIcon,
  },
  { id: "products" as EntityType, name: "Products", icon: PackageIcon },
  { id: "customers" as EntityType, name: "Customers", icon: UsersIcon },
  { id: "deliveries" as EntityType, name: "Deliveries", icon: TruckIcon },
  {
    id: "billing-documents" as EntityType,
    name: "Billing Documents",
    icon: FileTextIcon,
  },
  { id: "payments" as EntityType, name: "Payments", icon: CreditCardIcon },
]

function DataExplorer() {
  const [activeEntity, setActiveEntity] = useState<EntityType>("sales-orders")

  const { data, isLoading, error } = useQuery({
    queryKey: ["data", activeEntity],
    queryFn: () =>
      apiFetch<PaginatedResponse<Record<string, unknown>>>(
        `/api/v1/${activeEntity}`
      ),
  })

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="font-heading text-2xl font-bold tracking-tight">
          Data Explorer
        </h1>
        <p className="text-muted-foreground">
          Browse SAP order-to-cash data across different entity types.
        </p>
      </div>

      <div className="flex flex-wrap gap-2">
        {entityTypes.map((entity) => (
          <Button
            key={entity.id}
            variant={activeEntity === entity.id ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveEntity(entity.id)}
            className="gap-2"
          >
            <entity.icon className="size-4" />
            {entity.name}
          </Button>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            {(() => {
              const et = entityTypes.find((e) => e.id === activeEntity)
              const Icon = et?.icon
              return Icon ? <Icon className="size-5 text-primary" /> : null
            })()}
            {entityTypes.find((e) => e.id === activeEntity)?.name}
          </CardTitle>
          <CardDescription>
            {data?.total
              ? `Showing ${data.data.length} of ${data.total} records`
              : "No data available"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              Loading...
            </div>
          )}
          {error && (
            <div className="flex items-center justify-center py-8 text-destructive">
              Error loading data
            </div>
          )}
          {!isLoading && !error && data?.data && data.data.length > 0 && (
            <DataTable data={data.data} activeEntity={activeEntity} />
          )}
          {!isLoading && !error && (!data?.data || data.data.length === 0) && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              No data available
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

function DataTable({
  data,
  activeEntity,
}: {
  data: Array<Record<string, unknown>>
  activeEntity: EntityType
}) {
  const columns = getColumns(activeEntity)

  return (
    <Table className="text-sm">
      <TableHeader>
        <TableRow>
          {columns.map((col) => (
            <TableHead key={col.key} className="px-4 py-3">
              {col.label}
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((row, i) => (
          <TableRow key={i}>
            {columns.map((col) => (
              <TableCell key={col.key} className="px-4 py-3">
                <CellValue value={row[col.key]} />
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

function CellValue({ value }: { value: unknown }) {
  if (value == null) return <span className="text-muted-foreground">—</span>
  if (typeof value === "object") return <span>{JSON.stringify(value)}</span>
  return <span>{String(value)}</span>
}

interface ColumnDef {
  key: string
  label: string
}

function getColumns(entity: EntityType): Array<ColumnDef> {
  switch (entity) {
    case "sales-orders":
      return [
        { key: "sales_order", label: "Order #" },
        { key: "sold_to_party", label: "Sold-To Party" },
        { key: "creation_date", label: "Created" },
        { key: "total_net_amount", label: "Amount" },
        { key: "transaction_currency", label: "Currency" },
        { key: "overall_delivery_status", label: "Delivery Status" },
      ]
    case "products":
      return [
        { key: "product", label: "Product" },
        { key: "product_type", label: "Type" },
        { key: "product_group", label: "Group" },
        { key: "base_unit", label: "Unit" },
        { key: "division", label: "Division" },
      ]
    case "customers":
      return [
        { key: "business_partner", label: "Partner ID" },
        { key: "business_partner_full_name", label: "Name" },
        { key: "business_partner_category", label: "Category" },
        { key: "created_by_user", label: "Created By" },
      ]
    case "deliveries":
      return [
        { key: "delivery_document", label: "Delivery #" },
        { key: "shipping_point", label: "Shipping Point" },
        { key: "creation_date", label: "Created" },
        { key: "overall_goods_movement_status", label: "Goods Movement" },
        { key: "overall_picking_status", label: "Picking Status" },
      ]
    case "billing-documents":
      return [
        { key: "billing_document", label: "Billing #" },
        { key: "sold_to_party", label: "Sold-To Party" },
        { key: "billing_document_date", label: "Date" },
        { key: "total_net_amount", label: "Amount" },
        { key: "transaction_currency", label: "Currency" },
      ]
    case "payments":
      return [
        { key: "accounting_document", label: "Document #" },
        { key: "customer", label: "Customer" },
        { key: "amount_in_transaction_currency", label: "Amount" },
        { key: "transaction_currency", label: "Currency" },
        { key: "posting_date", label: "Posted" },
      ]
    default:
      return []
  }
}

export const Route = createFileRoute("/data/")({
  component: DataExplorer,
})
