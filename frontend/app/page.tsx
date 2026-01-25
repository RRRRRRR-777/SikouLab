/**
 * Home - Hello Worldページ
 *
 * バックエンドAPIからHello Worldメッセージを取得して表示する。
 */
async function getHelloMessage(): Promise<{ message: string }> {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const res = await fetch(apiUrl, {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Backend API error: ${res.status}`);
  }

  return res.json();
}

export default async function Home() {
  const data = await getHelloMessage();

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-center gap-8 px-16 bg-white dark:bg-black">
        <h1 className="text-4xl font-bold text-black dark:text-zinc-50">{data.message}</h1>
        <p className="text-lg text-zinc-600 dark:text-zinc-400">Frontend: http://localhost:3000</p>
        <p className="text-lg text-zinc-600 dark:text-zinc-400">Backend: http://localhost:8080</p>
      </main>
    </div>
  );
}
