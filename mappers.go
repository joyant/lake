package lake

import (
    "bytes"
    "encoding/xml"
    "io"
    "io/ioutil"
)

func Load(xmlFilename string) (err error) {
    var b []byte
    b, err = ioutil.ReadFile(xmlFilename)
    if err != nil {
        return
    }
    var (
        t xml.Token
        decoder = xml.NewDecoder(bytes.NewReader(b))
    )
    for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
        case xml.StartElement:
            if token.Name.Local == "mapper" {
                if len(token.Attr) > 0 && token.Attr[0].Name.Local == "resource" {
                    resource := token.Attr[0].Value
                    err = loadXML(resource)
                    if err != nil {
                        return
                    }
                    logger.Infof("loaded file %s success\n", resource)
                }
            }
        }
    }
    if err != io.EOF {
        return
    }
    return nil
}
