/**
 * Bridges a TypeScript object with a read() method to a Go io.Reader.
 * 
 * @param reader An object with a `read(buffer: ArrayBuffer): number` method.
 * @returns An opaque handle representing a Go io.Reader.
 */
declare function wrapReader(reader: { read(buffer: ArrayBuffer): number }): any;

/**
 * Bridges a TypeScript object with a write() method to a Go io.Writer.
 * 
 * @param writer An object with a `write(buffer: ArrayBuffer): number` method.
 * @returns An opaque handle representing a Go io.Writer.
 */
declare function wrapWriter(writer: { write(buffer: ArrayBuffer): number }): any;
