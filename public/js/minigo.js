(function () {
   var size = 9;

   var GoSquare = React.createClass({
       render: function() {
           return (
             <div className="square">
               <div className="stone empty"></div>
               <div className="grid"></div>
             </div>
           );
       }
   });

   var GoBoard = React.createClass({
       render: function() {
         var rows = [];
         for (var r = 0; r < size; r++) {
           var row = [];
           for (var c = 0; c < size; c++) {
             row.push(<GoSquare key={c} />);
           }
           rows.push(<div className="row" key={r}>{row}</div>);
         }
         return (
           <div className="goboard">
             {rows}
           </div>
         );
       }
   });

   ReactDOM.render(
     <GoBoard />,
     document.getElementById('content')
   );
 })();
