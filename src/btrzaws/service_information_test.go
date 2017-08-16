package btrzaws
import(
  "testing"
  "os"
  "encoding/json"
  "io/ioutil"
)
const(
  ServicesFile="../../sample_files/services.json"
)

func TestServiceAPIKeys(t*testing.T)  {
  fileData,err:=ioutil.ReadFile(ServicesFile)
  if err!=nil{
    t.Fatal(err)
  }
  info:=&ServiceInformation{}
  json.Unmarshal(fileData, info)

}
