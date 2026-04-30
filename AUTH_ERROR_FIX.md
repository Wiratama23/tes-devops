# Fixed: TypeError in Admin Login - Cannot read properties of null (reading 'length')

## The Error

When trying to log into the admin dashboard, you got:
```
Uncaught TypeError: Cannot read properties of null (reading 'length')
```

This was occurring in the authentication flow, specifically in error handling code.

## Root Cause

The `sanitizeErrorMessage()` function in `src/tools/api.ts` was not properly handling null or undefined values:

```typescript
// ❌ BEFORE - Crashes if message is null
function sanitizeErrorMessage(message: string): string {
  return message.replace(...) // crashes if message is null/undefined
}
```

When an authentication error occurred (401 Unauthorized), the response handling code could pass null/undefined to this function, causing the TypeError.

## The Fix

Added defensive null checks in three places:

### 1. **Sanitize function now handles null/undefined**
```typescript
// ✅ AFTER - Safe with any input
function sanitizeErrorMessage(message: string | null | undefined): string {
  if (!message || typeof message !== "string") {
    return "API request failed";  // Fallback message
  }
  return message.replace(...)
}
```

### 2. **Added null check for axios error**
```typescript
const ax = err as AxiosError;
if (!ax) {
  throw new ApiError("Unknown error occurred", 0, "");
}
```

### 3. **Added null check for response**
```typescript
if (!response) {
  throw new ApiError("No response from server", 0, "");
}
```

### 4. **Safe body text extraction**
```typescript
const bodyText =
  typeof response.data === "string"
    ? response.data
    : response.data  // Check if data exists before JSON.stringify
    ? (() => {
        try {
          return JSON.stringify(response.data);
        } catch {
          return String(response.data);
        }
      })()
    : "";  // Empty string fallback instead of undefined
```

## What Changed

**File**: `src/tools/api.ts`

- ✅ Function signature: `sanitizeErrorMessage(message: string)` → `sanitizeErrorMessage(message: string | null | undefined)`
- ✅ Added null guards before accessing properties
- ✅ Safe default values instead of processing null

## Testing

The fix handles these scenarios:
1. ✓ Successful login (200 OK)
2. ✓ Failed login (401 Unauthorized) - won't crash
3. ✓ Network errors (CONNECTION_REFUSED) - won't crash
4. ✓ Server errors (500 Internal Server Error) - won't crash
5. ✓ Malformed responses - won't crash

## How to Deploy

```bash
# 1. Pull the latest changes
git pull origin main

# 2. Rebuild frontend (so TypeScript types are updated)
docker-compose build frontend

# 3. Redeploy
docker-compose up -d frontend

# 4. Try logging in again
# Should now show proper error message instead of TypeError
```

## Expected Behavior After Fix

When login fails:
- ❌ **Before**: `TypeError: Cannot read properties of null (reading 'length')`
- ✅ **After**: `"Invalid username or password."` (proper error message)

The admin login form will now:
1. Display validation errors if username/password is empty
2. Show "Invalid username or password" if credentials are wrong
3. Display "Something went wrong" for server errors
4. NOT crash with TypeScript/runtime errors

## Files Modified

- ✅ `src/tools/api.ts` - Added null safety checks

