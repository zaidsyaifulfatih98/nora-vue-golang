// Opens a bare popup window with just the image, sized to the page, and
// triggers the browser print dialog once it's loaded — used to print a
// photobooth result without the app's own chrome around it.
export function usePrintImage() {
  function printImage(url: string) {
    const printWindow = window.open('', '_blank', 'width=800,height=1000')
    if (!printWindow) return

    printWindow.document.write(`
      <html>
        <head>
          <title>Print</title>
          <style>
            @page { margin: 0; }
            html, body { margin: 0; padding: 0; height: 100%; display: flex; align-items: center; justify-content: center; }
            img { max-width: 100%; max-height: 100vh; }
          </style>
        </head>
        <body>
          <img src="${url}" />
        </body>
      </html>
    `)
    printWindow.document.close()

    const image = printWindow.document.querySelector('img')
    const triggerPrint = () => {
      printWindow.focus()
      printWindow.print()
    }
    if (image?.complete) triggerPrint()
    else image?.addEventListener('load', triggerPrint)
  }

  return { printImage }
}
